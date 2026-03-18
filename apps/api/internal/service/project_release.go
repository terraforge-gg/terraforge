package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/terraforge-gg/terraforge/internal/database"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/lib/aws"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

type CreateProjectReleaseDependencyParams struct {
	ReleaseId        string
	ProjectId        string
	MinVersionNumber *string
	Type             string
}

type CreateReleaseParams struct {
	Name            string
	VersionNumber   string
	Changelog       *string
	LoaderVersionId string
	FileUrl         string
	Dependencies    []CreateProjectReleaseDependencyParams
}

type ProjectReleaseService interface {
	GetReleaseByIdWithDependencies(ctx context.Context, projectIdentifier string, releaseId string, userId string) (*models.ProjectRelease, error)
	CreateRelease(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error)
	GenerateProjectReleasePresignedPutUrl(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error)
	GetReleasesByProjectId(ctx context.Context, id string, userId string) ([]models.ProjectRelease, error)
}

type projectReleaseService struct {
	logger             *slog.Logger
	cdnUrl             string
	db                 *sql.DB
	projectRepo        repository.ProjectRepository
	projectReleaseRepo repository.ProjectReleaseRepository
	loaderVersionRepo  repository.LoaderVersionRepository
	objectStoreService ObjectStoreService
}

func NewProjectReleaseService(
	logger *slog.Logger,
	cdnUrl string,
	db *sql.DB,
	projectRepo repository.ProjectRepository,
	projectReleaseRepo repository.ProjectReleaseRepository,
	loaderVersionRepo repository.LoaderVersionRepository,
	objectStoreService ObjectStoreService) ProjectReleaseService {
	return &projectReleaseService{
		logger:             logger,
		cdnUrl:             cdnUrl,
		db:                 db,
		projectRepo:        projectRepo,
		projectReleaseRepo: projectReleaseRepo,
		loaderVersionRepo:  loaderVersionRepo,
		objectStoreService: objectStoreService,
	}
}

func (s *projectReleaseService) GetReleaseByIdWithDependencies(ctx context.Context, projectIdentifier string, releaseId string, userId string) (*models.ProjectRelease, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, projectIdentifier, userId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	release, err := s.projectReleaseRepo.FindReleaseByIdWithDependencies(ctx, s.db, releaseId)

	if err != nil {
		return nil, err
	}

	if release == nil {
		return nil, custom_errors.ErrProjectReleaseNotFound
	}

	return release, nil
}

func (s *projectReleaseService) CreateRelease(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, projectIdentifier, userId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, tx, project.Id, userId)

	if err != nil {
		return nil, err
	}

	if projectMember == nil {
		return nil, custom_errors.ErrProjectUnauthorisedAction
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return nil, custom_errors.ErrProjectUnauthorisedAction
	}

	loaderVersion, err := s.loaderVersionRepo.FindLoaderVersionById(ctx, tx, params.LoaderVersionId)

	if err != nil {
		return nil, err
	}

	if loaderVersion == nil {
		return nil, custom_errors.ErrLoaderVersionNotFound
	}

	parsedFileUrl, err := url.Parse(params.FileUrl)
	sourceKey := aws.ExtractS3Key(parsedFileUrl.Path)

	if err != nil {
		return nil, custom_errors.ErrProjectReleaseFailedToParseFileUrl
	}

	metadata, err := s.objectStoreService.GetFileMetadate(ctx, sourceKey)

	if err != nil {
		return nil, custom_errors.ErrProjectReleaseUploadedFileNotFound
	}

	release := models.ProjectRelease{
		Id:              utils.NewUUID(),
		ProjectId:       project.Id,
		Name:            params.Name,
		Changelog:       params.Changelog,
		VersionNumber:   params.VersionNumber,
		LoaderVersionId: loaderVersion.Id,
		LoaderVersion: models.LoaderVersion{
			Id:           loaderVersion.Id,
			GameVersion:  loaderVersion.GameVersion,
			VersionLabel: loaderVersion.VersionLabel,
			BuildType:    loaderVersion.BuildType,
			ReleasedAt:   loaderVersion.ReleasedAt,
			UpdatedAt:    loaderVersion.UpdatedAt,
		},
		Downloads: 0,
		FileUrl:   params.FileUrl,
		FileSize:  metadata.ContentLength,
		FileHash:  metadata.ETag,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	destinationKey := fmt.Sprintf("users/%s/projects/%s/releases/%s/%s_%s.tmod", userId, project.Id, release.Id, project.Slug, release.VersionNumber)

	newPath, err := s.objectStoreService.MoveFile(ctx, sourceKey, destinationKey)

	if err != nil {
		s.logger.Error("Failed to move file.", "Source", sourceKey, "Destination:", destinationKey)
		return nil, err
	}

	release.FileUrl = s.cdnUrl + newPath

	err = s.projectReleaseRepo.InsertRelease(ctx, tx, &release)

	if err != nil {
		switch {
		case errors.Is(err, database.ErrUniqueViolation):
			return nil, custom_errors.ErrProjectReleaseNumberAlreadyExists
		}
		return nil, err
	}

	seen := make(map[string]bool)

	for _, d := range params.Dependencies {
		key := d.ReleaseId

		if !seen[key] {
			seen[key] = true
		} else {
			return nil, custom_errors.ErrDuplicateProjectReleaseDependency
		}
	}

	deps := make([]models.ProjectReleaseDependency, 0, len(params.Dependencies))

	for _, dep := range params.Dependencies {
		p, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, dep.ProjectId, userId)

		if err != nil {
			return nil, err
		}

		if p == nil {
			return nil, custom_errors.ErrProjectReleaseDependencyNotFound
		}

		if p.Id == project.Id {
			return nil, custom_errors.ErrCircularProjectReleaseDependency
		}

		if dep.MinVersionNumber != nil {
			min, err := s.projectReleaseRepo.FindByProjectIdAndVersionNumber(ctx, tx, dep.ProjectId, *dep.MinVersionNumber)

			if err != nil {
				return nil, err
			}

			if min == nil {
				return nil, custom_errors.ErrProjectReleaseDependencyMinVersionDoesNotExist
			}
		}

		releaseDep := models.ProjectReleaseDependency{
			Id:                  utils.NewUUID(),
			ReleaseId:           release.Id,
			DependencyProjectId: dep.ProjectId,
			MinVersionNumber:    dep.MinVersionNumber,
			Type:                models.ProjectReleaseDependencyType(dep.Type),
			CreatedAt:           time.Now().UTC(),
		}

		deps = append(deps, releaseDep)
	}

	if len(deps) > 0 {
		err := s.projectReleaseRepo.InsertDependencies(ctx, tx, deps)

		if err != nil {
			return nil, err
		}
	}

	s.projectRepo.UpdateProject(ctx, tx, *project)

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	release.Dependencies = deps

	return &release, nil
}

func (s *projectReleaseService) GenerateProjectReleasePresignedPutUrl(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error) {
	fileSizeBytes, err := strconv.ParseInt(fileSize, 10, 64)

	if err != nil {
		return "", custom_errors.ErrProjectReleaseInvalidFileSize
	}

	if fileSizeBytes < 0 || fileSizeBytes > 524_288_000 {
		return "", custom_errors.ErrProjectReleaseInvalidFileSize
	}

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, projectIdentifier, userId)

	if err != nil {
		return "", err
	}

	if project == nil {
		return "", custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, s.db, project.Id, userId)

	if err != nil {
		return "", err
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return "", custom_errors.ErrProjectUnauthorisedAction
	}

	uploadId := utils.NewUUID()
	key := fmt.Sprintf("uploads/temp/%s/%s_%s.tmod", userId, uploadId, time.Now().UTC().Format(time.RFC3339))

	url, err := s.objectStoreService.GeneratePresignedPutUrl(ctx, key, "application/octet-stream", fileSizeBytes)

	if err != nil {
		s.logger.Error("Failed to generate project release presigned url.", "Error:", err)
		return "", err
	}

	return url, err
}

func (s *projectReleaseService) GetReleasesByProjectId(ctx context.Context, projectIdentifier string, userId string) ([]models.ProjectRelease, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, projectIdentifier, userId)

	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	releases, err := s.projectReleaseRepo.FindReleasesByProjectIdWithLoaderVersion(ctx, s.db, project.Id)

	if err != nil {
		return nil, err
	}

	return releases, nil
}
