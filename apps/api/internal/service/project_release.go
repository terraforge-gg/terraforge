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
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

type CreateProjectReleaseDependencyParams struct {
	VersionId        string
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
	GetReleaseById(ctx context.Context, id string) (*models.ProjectRelease, error)
	GetReleaseByIdWithDependencies(ctx context.Context, id string) (*models.ProjectRelease, error)
	CreateRelease(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error)
	GenerateProjectReleasePresignedPutUrl(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error)
	GetReleasesByProjectId(ctx context.Context, id string) ([]models.ProjectRelease, error)
}

type projectVersionService struct {
	logger             *slog.Logger
	cdnUrl             string
	db                 *sql.DB
	projectRepo        repository.ProjectRepository
	projectVersionRepo repository.ProjectReleaseRepository
	loaderVersionRepo  repository.LoaderVersionRepository
	objectStoreService ObjectStoreService
}

func NewProjectReleaseService(
	logger *slog.Logger,
	cdnUrl string,
	db *sql.DB,
	projectRepo repository.ProjectRepository,
	projectVersionRepo repository.ProjectReleaseRepository,
	loaderVersionRepo repository.LoaderVersionRepository,
	objectStoreService ObjectStoreService) ProjectReleaseService {
	return &projectVersionService{
		logger:             logger,
		cdnUrl:             cdnUrl,
		db:                 db,
		projectRepo:        projectRepo,
		projectVersionRepo: projectVersionRepo,
		loaderVersionRepo:  loaderVersionRepo,
		objectStoreService: objectStoreService,
	}
}

func (s *projectVersionService) GetReleaseById(ctx context.Context, id string) (*models.ProjectRelease, error) {
	version, err := s.projectVersionRepo.FindReleaseById(ctx, s.db, id)

	if err != nil {
		panic(err)
	}

	if version == nil {
		return nil, custom_errors.ErrProjectReleaseNotFound
	}

	return version, nil
}

func (s *projectVersionService) GetReleaseByIdWithDependencies(ctx context.Context, id string) (*models.ProjectRelease, error) {
	version, err := s.projectVersionRepo.FindReleaseByIdWithDependencies(ctx, s.db, id)

	if err != nil {
		panic(err)
	}

	if version == nil {
		return nil, custom_errors.ErrProjectReleaseNotFound
	}

	return version, nil
}

func (s *projectVersionService) CreateRelease(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error) {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, projectIdentifier, nil)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, tx, project.Id, userId)

	if err != nil {
		panic(err)
	}

	if projectMember == nil {
		return nil, custom_errors.ErrProjectUnauthorisedAction
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return nil, custom_errors.ErrProjectUnauthorisedAction
	}

	loaderVersion, err := s.loaderVersionRepo.FindLoaderVersionById(ctx, tx, params.LoaderVersionId)

	if err != nil {
		panic(err)
	}

	if loaderVersion == nil {
		return nil, custom_errors.ErrLoaderVersionNotFound
	}

	parsedFileUrl, err := url.Parse(params.FileUrl)
	sourceKey := utils.ExtractS3Key(parsedFileUrl.Path)

	if err != nil {
		return nil, custom_errors.ErrProjectReleaseFailedToParseFileUrl
	}

	metadata, err := s.objectStoreService.GetFileMetadate(ctx, sourceKey)

	if err != nil {
		return nil, custom_errors.ErrProjectReleaseUploadedFileNotFound
	}

	release := models.ProjectRelease{
		Id:              utils.GenerateUUID(),
		ProjectId:       project.Id,
		Name:            params.Name,
		Changelog:       params.Changelog,
		VersionNumber:   params.VersionNumber,
		LoaderVersionId: loaderVersion.Id,
		LoaderVersion: models.LoaderVersion{
			Id:              loaderVersion.Id,
			GameVersion:     loaderVersion.GameVersion,
			InternalVersion: loaderVersion.InternalVersion,
			Status:          loaderVersion.Status,
			IsLegacy:        loaderVersion.IsLegacy,
			ReleasedAt:      loaderVersion.ReleasedAt,
		},
		Downloads: 0,
		FileUrl:   params.FileUrl,
		FileSize:  metadata.ContentLength,
		FileHash:  metadata.ETag,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	destinationKey := fmt.Sprintf("users/%s/projects/%s/versions/%s/%s_%s.tmod", userId, project.Id, release.Id, project.Slug, release.VersionNumber)

	newPath, err := s.objectStoreService.MoveFile(ctx, sourceKey, destinationKey)

	if err != nil {
		s.logger.Error("Failed to move file.", "Source", sourceKey, "Destination:", destinationKey)
		panic(err)
	}

	release.FileUrl = s.cdnUrl + newPath

	err = s.projectVersionRepo.InsertRelease(ctx, tx, &release)

	if err != nil {
		switch {
		case errors.Is(err, database.ErrUniqueViolation):
			return nil, custom_errors.ErrProjectReleaseNumberAlreadyExists
		}
		panic(err)
	}

	seen := make(map[string]bool)

	for _, d := range params.Dependencies {
		key := d.VersionId

		if !seen[key] {
			seen[key] = true
		} else {
			return nil, custom_errors.ErrDuplicateProjectReleaseDependency
		}
	}

	deps := make([]models.ProjectReleaseDependency, 0, len(params.Dependencies))

	for _, dep := range params.Dependencies {
		p, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, dep.ProjectId, nil)

		if err != nil {
			panic(err)
		}

		if p == nil {
			return nil, custom_errors.ErrProjectReleaseDependencyNotFound
		}

		if p.Id == project.Id {
			return nil, custom_errors.ErrCircularProjectReleaseDependency
		}

		if dep.MinVersionNumber != nil {
			min, err := s.projectVersionRepo.FindByProjectIdAndVersionNumber(ctx, tx, dep.ProjectId, *dep.MinVersionNumber)

			if err != nil {
				panic(err)
			}

			if min == nil {
				return nil, custom_errors.ErrProjectReleaseDependencyMinVersionDoesNotExist
			}
		}

		releaseDep := models.ProjectReleaseDependency{
			Id:                  utils.GenerateUUID(),
			ReleaseId:           release.Id,
			DependencyProjectId: dep.ProjectId,
			MinVersionNumber:    dep.MinVersionNumber,
			Type:                models.ProjectReleaseDependencyType(dep.Type),
			CreatedAt:           time.Now().UTC(),
		}

		deps = append(deps, releaseDep)
	}

	if len(deps) > 0 {
		err := s.projectVersionRepo.InsertDependencies(ctx, tx, deps)

		if err != nil {
			panic(err)
		}
	}

	s.projectRepo.UpdateProject(ctx, tx, *project)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

	release.Dependencies = deps

	return &release, nil
}

func (s *projectVersionService) GenerateProjectReleasePresignedPutUrl(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error) {
	fileSizeBytes, err := strconv.ParseInt(fileSize, 10, 64)

	if err != nil {
		return "", custom_errors.ErrProjectReleaseInvalidFileSize
	}

	if fileSizeBytes < 0 || fileSizeBytes > 524_288_000 {
		return "", custom_errors.ErrProjectReleaseInvalidFileSize
	}

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, projectIdentifier, nil)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return "", custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, s.db, project.Id, userId)

	if err != nil {
		panic(err)
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return "", custom_errors.ErrProjectUnauthorisedAction
	}

	uploadId := utils.GenerateUUID()
	key := fmt.Sprintf("uploads/temp/%s/%s_%d.tmod", userId, uploadId, time.Now().Unix())

	url, err := s.objectStoreService.GeneratePresignedPutUrl(ctx, key, "application/octet-stream", fileSizeBytes)

	if err != nil {
		s.logger.Error("Failed to generate project version presigned url.", "Error:", err)
		panic(err)
	}

	return url, err
}

func (s *projectVersionService) GetReleasesByProjectId(ctx context.Context, projectIdentifier string) ([]models.ProjectRelease, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, projectIdentifier, nil)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	versions, err := s.projectVersionRepo.FindReleasesByProjectIdWithLoaderVersion(ctx, s.db, project.Id)

	if err != nil {
		panic(err)
	}

	return versions, nil
}
