package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
)

type LoaderVersionService interface {
	GetLoaderVersionById(ctx context.Context, id string) (*models.LoaderVersion, error)
	GetLoaderVersionByGameVersion(ctx context.Context, gameVersion string) (*models.LoaderVersion, error)
	GetLoaderVersionByLabel(ctx context.Context, label string) (*models.LoaderVersion, error)
	GetLoaderVersions(ctx context.Context) ([]models.LoaderVersion, error)
	CreateLoaderVersion(ctx context.Context, params CreateLoaderVersionParams) error
}

type loaderVersionService struct {
	logger            *slog.Logger
	db                *sql.DB
	loaderVersionRepo repository.LoaderVersionRepository
}

func NewLoaderVersionService(logger *slog.Logger, db *sql.DB, loaderVersionRepo repository.LoaderVersionRepository) LoaderVersionService {
	return &loaderVersionService{logger: logger, db: db, loaderVersionRepo: loaderVersionRepo}
}

func (s *loaderVersionService) GetLoaderVersionById(ctx context.Context, id string) (*models.LoaderVersion, error) {
	modLoaderVersion, err := s.loaderVersionRepo.FindLoaderVersionById(ctx, s.db, id)

	if err != nil {
		return nil, err
	}

	if modLoaderVersion == nil {
		return nil, errors.ErrLoaderVersionNotFound
	}

	return modLoaderVersion, nil
}

func (s *loaderVersionService) GetLoaderVersionByGameVersion(ctx context.Context, gameVersion string) (*models.LoaderVersion, error) {
	loaderVersion, err := s.loaderVersionRepo.FindLoaderVersionByGameVersion(ctx, s.db, gameVersion)

	if err != nil {
		return nil, err
	}

	if loaderVersion == nil {
		return nil, errors.ErrLoaderVersionNotFound
	}

	return loaderVersion, nil
}

func (s *loaderVersionService) GetLoaderVersionByLabel(ctx context.Context, label string) (*models.LoaderVersion, error) {
	loaderVersion, err := s.loaderVersionRepo.FindLoaderVersionByLabel(ctx, s.db, label)

	if err != nil {
		return nil, err
	}

	if loaderVersion == nil {
		return nil, errors.ErrLoaderVersionNotFound
	}

	return loaderVersion, nil
}

func (s *loaderVersionService) GetLoaderVersions(ctx context.Context) ([]models.LoaderVersion, error) {
	loaderVersions, err := s.loaderVersionRepo.FindLoaderVersions(ctx, s.db)

	if err != nil {
		return nil, err
	}

	return loaderVersions, nil
}

type CreateLoaderVersionParams struct {
	Id           string
	GameVersion  string
	VersionLabel string
	BuildType    string
	ReleasedAt   time.Time
	UpdatedAt    time.Time
}

func (s *loaderVersionService) CreateLoaderVersion(ctx context.Context, params CreateLoaderVersionParams) error {
	err := s.loaderVersionRepo.InsertLoaderVersion(ctx, s.db, &models.LoaderVersion{
		Id:           params.Id,
		GameVersion:  params.GameVersion,
		VersionLabel: params.VersionLabel,
		BuildType:    models.LoaderVersionBuildType(params.BuildType),
		ReleasedAt:   params.ReleasedAt,
		UpdatedAt:    params.UpdatedAt,
	})

	if err != nil {
		return err
	}

	return nil
}
