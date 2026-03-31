package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type ProjectRepository interface {
	InsertProject(ctx context.Context, q database.Querier, project *models.Project) error
	InsertProjectMember(ctx context.Context, q database.Querier, projectMember *models.ProjectMember) error
	FindProjectByIdentifier(ctx context.Context, q database.Querier, projectIdentifier string, userId string) (*models.Project, error)
	FindProjectMembersByProjectIdentifier(ctx context.Context, q database.Querier, projectIdentifier string, userId string) ([]models.ProjectMember, error)
	FindProjectMemberByProjectIdAndUserId(ctx context.Context, q database.Querier, projectId string, userId string) (*models.ProjectMember, error)
	UpdateProject(ctx context.Context, q database.Querier, project models.Project) error
	DeleteProjectByIdentifier(ctx context.Context, q database.Querier, identifier string, deletedAt time.Time) error
	FindProjectsByUserIdentifier(ctx context.Context, q database.Querier, userIdentifier string, status models.ProjectStatus) ([]models.Project, error)
}

type projectRepository struct{}

func NewProjectRepository() ProjectRepository {
	return &projectRepository{}
}

func (r *projectRepository) InsertProject(ctx context.Context, q database.Querier, project *models.Project) error {
	query := `INSERT INTO "project"
		("id", "name", "slug", "summary", "description", "iconUrl", "downloads", "type", "status", "createdAt", "updatedAt", "userId")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	_, err := q.ExecContext(
		ctx,
		query,
		project.Id,
		project.Name,
		project.Slug,
		project.Summary,
		project.Description,
		project.IconUrl,
		project.Downloads,
		project.Type,
		project.Status,
		project.CreatedAt,
		project.UpdatedAt,
		project.UserId)

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

func (r *projectRepository) FindProjectByIdentifier(ctx context.Context, q database.Querier, projectIdentifier string, userId string) (*models.Project, error) {
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
		FROM "active_project" 
		WHERE ("id" = $1 OR "slug" = $1)
			AND (
				"status" = 'approved'
				OR ("status" = 'draft' AND "userId" = $2)
			);`

	project := &models.Project{}
	err := q.QueryRowContext(ctx, query, projectIdentifier, userId).Scan(
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

func (r *projectRepository) FindProjectMembersByProjectIdentifier(ctx context.Context, q database.Querier, projectIdentifier string, userId string) ([]models.ProjectMember, error) {
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
        JOIN "active_project" p ON pm."projectId" = p."id"
        WHERE (pm."projectId" = $1 OR p."slug" = $1)
			AND (
				p."status" = 'approved'
				OR (p."status" = 'draft' AND EXISTS (
					SELECT 1 FROM "project_member" pm2
					WHERE pm2."projectId" = p."id"
					AND pm2."userId" = $2
				))
			);`

	rows, err := q.QueryContext(ctx, query, projectIdentifier, userId)

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

	if len(projectMembers) == 0 {
		return nil, nil
	}

	return projectMembers, nil
}

func (r *projectRepository) FindProjectMemberByProjectIdAndUserId(ctx context.Context, q database.Querier, projectId string, userId string) (*models.ProjectMember, error) {
	query := `
		SELECT
			"id",
			"projectId",
			"userId",
			"role",
			"createdAt"
		FROM "project_member" 
		WHERE ("projectId" = $1 AND "userId" = $2);`

	projectMember := &models.ProjectMember{}
	err := q.QueryRowContext(ctx, query, projectId, userId).Scan(
		&projectMember.Id,
		&projectMember.ProjectId,
		&projectMember.UserId,
		&projectMember.Role,
		&projectMember.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return projectMember, nil
}

func (r *projectRepository) UpdateProject(ctx context.Context, q database.Querier, project models.Project) error {
	query := `
        UPDATE "active_project"
        SET "name" = $2,
			"slug" = $3,
            "summary" = $4,
            "description" = $5,
			"iconUrl" = $6,
            "updatedAt" = now()
        WHERE id = $1;
    `

	_, err := q.ExecContext(ctx, query,
		project.Id,
		project.Name,
		project.Slug,
		project.Summary,
		project.Description,
		project.IconUrl)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectRepository) DeleteProjectByIdentifier(ctx context.Context, q database.Querier, identifier string, deletedAt time.Time) error {
	query := `
        UPDATE "active_project"
        SET "deletedAt" = $2
        WHERE "id" = $1 OR "slug" = $1;
	`
	_, err := q.ExecContext(ctx, query, identifier, deletedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *projectRepository) FindProjectsByUserIdentifier(ctx context.Context, q database.Querier, userIdentifier string, projectStatus models.ProjectStatus) ([]models.Project, error) {

	query := `
		SELECT
			p."id",
			p."name",
			p."slug",
			p."summary",
			p."description",
			p."iconUrl",
			p."downloads",
			p."type",
			p."status",
			p."createdAt",
			p."updatedAt",
			p."deletedAt",
			p."userId"
		FROM "active_project" p
		JOIN "user" u ON p."userId" = u."id"
		WHERE (u."id" = $1 OR u."username" = $1) AND p."status" = $2
		ORDER BY p."downloads" DESC, p."updatedAt" DESC LIMIT 100;
	`

	rows, err := q.QueryContext(ctx, query, userIdentifier, projectStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project

		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Slug,
			&p.Summary,
			&p.Description,
			&p.IconUrl,
			&p.Downloads,
			&p.Type,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
			&p.UserId,
		)
		if err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}
