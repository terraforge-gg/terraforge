package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type ProjectReleaseRepository interface {
	FindReleaseById(ctx context.Context, q database.Querier, id string) (*models.ProjectRelease, error)
	FindReleaseByIdWithDependencies(ctx context.Context, q database.Querier, id string) (*models.ProjectRelease, error)
	InsertRelease(ctx context.Context, q database.Querier, version *models.ProjectRelease) error
	InsertDependencies(ctx context.Context, q database.Querier, deps []models.ProjectReleaseDependency) error
	FindByProjectIdAndVersionNumber(ctx context.Context, q database.Querier, projectId string, versionNumber string) (*models.ProjectRelease, error)
	FindReleasesByProjectIdWithLoaderVersion(ctx context.Context, q database.Querier, projectId string) ([]models.ProjectRelease, error)
}

type projectReleaseRepository struct{}

func NewProjectReleaseRepository() ProjectReleaseRepository {
	return &projectReleaseRepository{}
}

func (r *projectReleaseRepository) FindReleaseById(ctx context.Context, q database.Querier, id string) (*models.ProjectRelease, error) {
	query := `
		SELECT
			"id",
			"projectId",
			"name",
			"changelog",
			"versionNumber",
			"loaderVersionId",
			"downloads",
			"fileUrl",
			"fileSize",
			"fileHash",
			"createdAt",
			"updatedAt",
			"publishedAt"
        FROM "project_release"
		WHERE id = $1;
	`

	version := &models.ProjectRelease{}

	err := q.QueryRowContext(ctx, query, id).Scan(
		&version.Id,
		&version.ProjectId,
		&version.Name,
		&version.Changelog,
		&version.VersionNumber,
		&version.LoaderVersionId,
		&version.Downloads,
		&version.FileUrl,
		&version.FileSize,
		&version.FileHash,
		&version.CreatedAt,
		&version.UpdatedAt,
		&version.PublishedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (r *projectReleaseRepository) FindReleaseByIdWithDependencies(ctx context.Context, q database.Querier, id string) (*models.ProjectRelease, error) {
	query := `
        SELECT
			v."id",
			v."projectId",
			v."name",
			v."changelog",
			v."versionNumber",
			v."loaderVersionId",
			v."downloads",
			v."fileUrl",
			v."fileSize",
			v."fileHash",
			v."createdAt",
			v."updatedAt",
			v."publishedAt",
			d."id",
			d."versionId",
			d."dependencyProjectId",
			d."minVersionNumber",
			d."type",
			d."createdAt",
			l."id",
			l."gameVersion",
			l."internalVersion",
			l."status",
			l."isLegacy",
			l."releasedAt",
			l."updatedAt"
		FROM "project_release" v
		LEFT JOIN "project_release_dependency" d ON v."id" = d."versionId"
		JOIN "loader_version" l ON v."loaderVersionId" = l."id"
		WHERE v."id" = $1;
    `

	rows, err := q.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var version *models.ProjectRelease
	var loaderVersion models.LoaderVersion
	var deps []models.ProjectReleaseDependency

	for rows.Next() {
		var (
			v                   models.ProjectRelease
			dep                 models.ProjectReleaseDependency
			depId               sql.NullString
			depReleaseId        sql.NullString
			depProjectId        string
			depMinVersionNumber sql.NullString
			depType             sql.NullString
			depCreatedAt        sql.NullTime
		)
		err := rows.Scan(
			&v.Id, &v.ProjectId, &v.Name, &v.Changelog, &v.VersionNumber, &v.LoaderVersionId,
			&v.Downloads, &v.FileUrl, &v.FileSize, &v.FileHash, &v.CreatedAt, &v.UpdatedAt, &v.PublishedAt,
			&depId, &depReleaseId, &depProjectId, &depMinVersionNumber, &depType, &depCreatedAt, &loaderVersion.Id,
			&loaderVersion.GameVersion,
			&loaderVersion.InternalVersion,
			&loaderVersion.Status,
			&loaderVersion.IsLegacy,
			&loaderVersion.ReleasedAt,
		)
		if err != nil {
			return nil, err
		}
		if version == nil {
			version = &v
		}
		if depId.Valid {
			dep.Id = depId.String
			dep.ReleaseId = depReleaseId.String
			dep.DependencyProjectId = depProjectId
			dep.MinVersionNumber = &depMinVersionNumber.String
			dep.Type = models.ProjectReleaseDependencyType(depType.String)
			dep.CreatedAt = depCreatedAt.Time
			deps = append(deps, dep)
		}
	}

	if version == nil {
		return nil, nil
	}

	version.LoaderVersion = loaderVersion
	version.Dependencies = deps

	return version, nil
}

func (r *projectReleaseRepository) InsertRelease(ctx context.Context, q database.Querier, version *models.ProjectRelease) error {
	query := `INSERT INTO "project_release" (
        "id", "projectId", "name", "changelog", "versionNumber", "loaderVersionId",
        "downloads", "fileUrl", "fileSize", "fileHash", "createdAt", "updatedAt", "publishedAt"
    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);`

	_, err := q.ExecContext(
		ctx,
		query,
		version.Id,
		version.ProjectId,
		version.Name,
		version.Changelog,
		version.VersionNumber,
		version.LoaderVersion.Id,
		version.Downloads,
		version.FileUrl,
		version.FileSize,
		version.FileHash,
		version.CreatedAt,
		version.UpdatedAt,
		version.PublishedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return database.ErrUniqueViolation
			}
		}

		return err
	}

	return nil
}

func (r *projectReleaseRepository) InsertDependencies(ctx context.Context, q database.Querier, deps []models.ProjectReleaseDependency) error {
	if len(deps) == 0 {
		return nil
	}

	valueParts := make([]string, 0, len(deps))
	args := make([]any, 0, len(deps)*5)

	for i, d := range deps {
		base := i * 5
		valueParts = append(valueParts, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", base+1, base+2, base+3, base+4, base+5))
		args = append(args, d.Id, d.ReleaseId, d.DependencyProjectId, d.Type, d.CreatedAt)
	}

	query := `INSERT INTO "project_release_dependency" ("id", "releaseId", "dependencyProjectId", "type", "createdAt") VALUES ` +
		strings.Join(valueParts, ",")

	_, err := q.ExecContext(ctx, query, args...)

	return err
}

func (r *projectReleaseRepository) FindByProjectIdAndVersionNumber(ctx context.Context, q database.Querier, projectId string, versionNumber string) (*models.ProjectRelease, error) {
	query := `
		SELECT
			"id",
			"projectId",
			"name",
			"changelog",
			"versionNumber",
			"loaderVersionId",
			"downloads",
			"fileUrl",
			"fileSize",
			"fileHash",
			"createdAt",
			"updatedAt",
			"publishedAt"
        FROM "project_release"
		WHERE "projectId" = $1 AND "versionNumber" = $2;
	`

	version := &models.ProjectRelease{}

	err := q.QueryRowContext(ctx, query, projectId, versionNumber).Scan(
		&version.Id,
		&version.ProjectId,
		&version.Name,
		&version.Changelog,
		&version.VersionNumber,
		&version.LoaderVersionId,
		&version.Downloads,
		&version.FileUrl,
		&version.FileSize,
		&version.FileHash,
		&version.CreatedAt,
		&version.UpdatedAt,
		&version.PublishedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (r *projectReleaseRepository) FindReleasesByProjectIdWithLoaderVersion(ctx context.Context, q database.Querier, projectId string) ([]models.ProjectRelease, error) {
	query := `
		SELECT
			v."id",
			v."projectId",
			v."name",
			v."changelog",
			v."versionNumber",
			v."loaderVersionId",
			v."downloads",
			v."fileUrl",
			v."fileSize",
			v."fileHash",
			v."createdAt",
			v."updatedAt",
			v."publishedAt",
			l."id",
			l."gameVersion",
			l."internalGameVersion",
			l."status",
			l."isLegacy",
			l."releasedAt",
			l."updatedAt"
		FROM "project_release" v
		LEFT JOIN "loader_version" l ON v."loaderVersionId" = l."id"
		WHERE v."projectId" = $1
		ORDER BY v."createdAt" DESC;`

	rows, err := q.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.ProjectRelease

	for rows.Next() {
		var v models.ProjectRelease

		var loaderId sql.NullString
		var loaderGameVersionStr sql.NullString
		var loaderInternalVersionStr sql.NullString
		var loaderStatus sql.NullString
		var loaderIsLegacy sql.NullBool
		var loaderReleasedAt time.Time
		var loaderUpdatedAt time.Time

		if err := rows.Scan(
			&v.Id,
			&v.ProjectId,
			&v.Name,
			&v.Changelog,
			&v.VersionNumber,
			&v.LoaderVersionId,
			&v.Downloads,
			&v.FileUrl,
			&v.FileSize,
			&v.FileHash,
			&v.CreatedAt,
			&v.UpdatedAt,
			&v.PublishedAt,
			&loaderId,
			&loaderGameVersionStr,
			&loaderInternalVersionStr,
			&loaderStatus,
			&loaderIsLegacy,
			&loaderReleasedAt,
			&loaderUpdatedAt,
		); err != nil {
			return nil, err
		}

		if loaderId.Valid {
			var lv models.LoaderVersion
			lv.Id = loaderId.String
			lv.GameVersion = loaderGameVersionStr.String
			lv.InternalVersion = loaderInternalVersionStr.String
			lv.Status = models.LoaderVersionStatus(loaderStatus.String)
			lv.IsLegacy = loaderIsLegacy.Bool
			lv.ReleasedAt = loaderReleasedAt
			lv.UpdatedAt = loaderUpdatedAt
			v.LoaderVersion = lv
		}

		versions = append(versions, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, nil
	}

	return versions, nil
}
