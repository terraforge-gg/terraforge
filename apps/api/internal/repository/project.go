package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type ProjectRepository interface {
	InsertProject(ctx context.Context, q database.Querier, project *models.Project) error
	InsertProjectMember(ctx context.Context, q database.Querier, projectMember *models.ProjectMember) error
	FindProjectByIdentifier(ctx context.Context, q database.Querier, identifier string, projectType models.ProjectType) (*models.Project, error)
	FindProjectByIdentifierIncludeDeleted(ctx context.Context, q database.Querier, identifier string, projectType models.ProjectType) (*models.Project, error)
	FindProjectMembersByProjectId(ctx context.Context, q database.Querier, projectId string) ([]models.ProjectMember, error)
}

type projectRepository struct{}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) InsertProject(ctx context.Context, q database.Querier, project *models.Project) error {
	query := `INSERT INTO "project"
		("id", "name", "slug", "summary", "type", "status", "createdAt", "updatedAt", "userId")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

	_, err := q.ExecContext(
		ctx,
		query,
		project.Id,
		project.Name,
		project.Slug,
		project.Summary,
		project.Type,
		project.Status,
		project.CreatedAt,
		project.UpdatedAt,
		project.UserId)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectRepository) InsertProjectMember(ctx context.Context, q database.Querier, projectMember *models.ProjectMember) error {
	query := `INSERT INTO "project_member"
		("id", "projectId", "userId", "role", "createdAt")
		VALUES ($1, $2, $3, $4, $5);`

	_, err := q.ExecContext(
		ctx,
		query,
		projectMember.Id,
		projectMember.ProjectId,
		projectMember.UserId,
		projectMember.Role,
		projectMember.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectRepository) FindProjectByIdentifier(ctx context.Context, q database.Querier, identifier string, projectType models.ProjectType) (*models.Project, error) {
	query := `
		SELECT
			"id",
			"name",
			"slug",
			"summary",
			"description",
			"iconUrl",
			"downloads",
			"type",
			"status",
			"createdAt",
			"updatedAt",
			"userId"
		FROM "project" 
		WHERE ("id" = $1 OR "slug" = $1) AND "type" IS $2 AND "deletedAt" IS NULL;`

	project := &models.Project{}
	err := q.QueryRowContext(ctx, query, identifier, projectType).Scan(
		&project.Id,
		&project.Name,
		&project.Slug,
		&project.Summary,
		&project.Description,
		&project.IconUrl,
		&project.Downloads,
		&project.Type,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *projectRepository) FindProjectByIdentifierIncludeDeleted(ctx context.Context, q database.Querier, identifier string, projectType models.ProjectType) (*models.Project, error) {
	query := `
		SELECT
			"id",
			"name",
			"slug",
			"summary",
			"description",
			"iconUrl",
			"downloads",
			"type",
			"status",
			"createdAt",
			"updatedAt",
			"userId"
		FROM "project" 
		WHERE ("id" = $1 OR "slug" = $1) AND "type" = $2;`

	project := &models.Project{}
	err := q.QueryRowContext(ctx, query, identifier, projectType).Scan(
		&project.Id,
		&project.Name,
		&project.Slug,
		&project.Summary,
		&project.Description,
		&project.IconUrl,
		&project.Downloads,
		&project.Type,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.UserId)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return project, nil
}

type projectMemberRow struct {
	Id        string
	ProjectId string
	UserId    string
	Role      string
	CreatedAt time.Time
	Username  string
	Image     sql.NullString
}

func (r *projectRepository) FindProjectMembersByProjectId(ctx context.Context, q database.Querier, projectId string) ([]models.ProjectMember, error) {
	query := `
        SELECT
            pm."id",
            pm."projectId",
            pm."userId",
            pm."role",
            pm."createdAt",
            u."username",
            u."image"
        FROM "project_member" pm
        JOIN "user" u ON pm."userId" = u."id"
        WHERE pm."projectId" = $1`

	rows, err := q.QueryContext(ctx, query, projectId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	projectMembers := []models.ProjectMember{}
	for rows.Next() {
		var row projectMemberRow
		if err := rows.Scan(
			&row.Id,
			&row.ProjectId,
			&row.UserId,
			&row.Role,
			&row.CreatedAt,
			&row.Username,
			&row.Image,
		); err != nil {
			return nil, err
		}

		var image *string
		if row.Image.Valid {
			image = &row.Image.String
		} else {
			image = nil
		}

		projectMembers = append(projectMembers, models.ProjectMember{
			Id:        row.Id,
			ProjectId: row.ProjectId,
			UserId:    row.UserId,
			Role:      models.ProjectMemberRole(row.Role),
			CreatedAt: row.CreatedAt,
			User: models.User{
				Username: row.Username,
				Image:    image,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projectMembers, nil
}
