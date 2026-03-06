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
	GetLoaderVersionByGameVersion(ctx context.Context, id string) (*models.LoaderVersion, error)
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
		panic(err)
	}

	if modLoaderVersion == nil {
		return nil, errors.ErrLoaderVersionNotFound
	}

	return modLoaderVersion, nil
}

func (s *loaderVersionService) GetLoaderVersionByGameVersion(ctx context.Context, gameVersion string) (*models.LoaderVersion, error) {
	loaderVersion, err := s.loaderVersionRepo.FindLoaderVersionByGameVersion(ctx, s.db, gameVersion)

	if err != nil {
		panic(err)
	}

	if loaderVersion == nil {
		return nil, errors.ErrLoaderVersionNotFound
	}

	return loaderVersion, nil
}

func (s *loaderVersionService) GetLoaderVersions(ctx context.Context) ([]models.LoaderVersion, error) {
	loaderVersions, err := s.loaderVersionRepo.FindLoaderVersions(ctx, s.db)

	if err != nil {
		panic(err)
	}

	return loaderVersions, nil
}

type CreateLoaderVersionParams struct {
	Id              string
	GameVersion     string
	InternalVersion string
	Status          string
	IsLegacy        bool
	ReleasedAt      time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s *loaderVersionService) CreateLoaderVersion(ctx context.Context, params CreateLoaderVersionParams) error {
	err := s.loaderVersionRepo.InsertLoaderVersion(ctx, s.db, &models.LoaderVersion{
		Id:              params.Id,
		GameVersion:     params.GameVersion,
		InternalVersion: params.InternalVersion,
		Status:          models.LoaderVersionStatus(params.Status),
		IsLegacy:        params.IsLegacy,
		ReleasedAt:      params.ReleasedAt,
		CreatedAt:       params.CreatedAt,
		UpdatedAt:       params.UpdatedAt,
	})

	if err != nil {
		panic(err)
	}

	return nil
}
