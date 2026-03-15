package repository

import (
	"context"
	"database/sql"

	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type LoaderVersionRepository interface {
	FindLoaderVersionById(ctx context.Context, q database.Querier, id string) (*models.LoaderVersion, error)
	FindLoaderVersionByGameVersion(ctx context.Context, q database.Querier, id string) (*models.LoaderVersion, error)
	FindLoaderVersions(ctx context.Context, q database.Querier) ([]models.LoaderVersion, error)
	InsertLoaderVersion(ctx context.Context, q database.Querier, loaderVersion *models.LoaderVersion) error
}

type loaderVersionRepository struct{}

func NewLoaderVersionRepository() LoaderVersionRepository {
	return &loaderVersionRepository{}
}

func (r *loaderVersionRepository) FindLoaderVersionById(ctx context.Context, q database.Querier, id string) (*models.LoaderVersion, error) {
	query := `
		SELECT 
			"id",
			"gameVersion",
			"versionLabel",
			"buildType",
			"releasedAt",
			"updatedAt"
		FROM "loader_version" 
		WHERE "id" = $1;
	`

	lv := &models.LoaderVersion{}

	err := q.QueryRow(query, id).Scan(
		&lv.Id,
		&lv.GameVersion,
		&lv.VersionLabel,
		&lv.BuildType,
		&lv.ReleasedAt,
		&lv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return lv, nil
}

func (r *loaderVersionRepository) FindLoaderVersionByGameVersion(ctx context.Context, q database.Querier, id string) (*models.LoaderVersion, error) {
	query := `
		SELECT 
			"id",
			"gameVersion",
			"versionLabel",
			"buildType",
			"releasedAt",
			"updatedAt"
		FROM "loader_version" 
		WHERE "gameVersion" = $1;
	`

	lv := &models.LoaderVersion{}

	err := q.QueryRow(query, id).Scan(
		&lv.Id,
		&lv.GameVersion,
		&lv.VersionLabel,
		&lv.BuildType,
		&lv.ReleasedAt,
		&lv.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return lv, nil
}

func (r *loaderVersionRepository) FindLoaderVersions(ctx context.Context, q database.Querier) ([]models.LoaderVersion, error) {
	query := `
		SELECT
			"id",
			"gameVersion",
			"versionLabel",
			"buildType",
			"releasedAt",
			"updatedAt"
		FROM "loader_version"
		ORDER BY "gameVersion" DESC, "buildType" ASC
		LIMIT 10;
	`

	rows, err := q.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var loaderVersions []models.LoaderVersion

	for rows.Next() {
		var lv models.LoaderVersion
		err := rows.Scan(
			&lv.Id,
			&lv.GameVersion,
			&lv.VersionLabel,
			&lv.BuildType,
			&lv.ReleasedAt,
			&lv.UpdatedAt)

		if err != nil {
			return nil, err
		}

		loaderVersions = append(loaderVersions, lv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return loaderVersions, nil
}

func (r *loaderVersionRepository) InsertLoaderVersion(ctx context.Context, q database.Querier, loaderVersion *models.LoaderVersion) error {
	query := `INSERT INTO "loader_version"
		("id", "gameVersion", "versionLabel", "buildType", "releasedAt", "updatedAt")
		VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := q.ExecContext(
		ctx,
		query,
		loaderVersion.Id,
		loaderVersion.GameVersion,
		loaderVersion.VersionLabel,
		loaderVersion.BuildType,
		loaderVersion.ReleasedAt,
		loaderVersion.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
