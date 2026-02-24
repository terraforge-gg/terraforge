package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

type ProjectService interface {
	CreateUserProject(ctx context.Context, params CreateUserProjectParams) (*models.Project, error)
	GetProjectByIdentifier(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error)
	GetProjectMembers(ctx context.Context, params GetProjectByIdentifierParams) ([]models.ProjectMember, error)
}

type projectService struct {
	logger      *slog.Logger
	db          *sql.DB
	projectRepo repository.ProjectRepository
}

func NewProjectService(logger *slog.Logger, db *sql.DB, projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{logger: logger, db: db, projectRepo: projectRepo}
}

type CreateUserProjectParams struct {
	Name    string
	Slug    string
	Summary *string
	Type    models.ProjectType
	UserId  string
}

func (s *projectService) CreateUserProject(ctx context.Context, params CreateUserProjectParams) (*models.Project, error) {
	project := &models.Project{
		Id:        utils.GenerateUUID(),
		Name:      params.Name,
		Slug:      params.Slug,
		Summary:   params.Summary,
		Type:      params.Type,
		Status:    models.ProjectStatusDraft,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserId:    params.UserId,
	}

	exists, err := s.projectRepo.FindProjectByIdentifierIncludeDeleted(ctx, s.db, project.Slug, project.Type)

	if err != nil {
		panic(err)
	}

	if exists != nil {
		return nil, errors.ErrProjectSlugUsed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err != nil {
		panic(err)
	}

	err = s.projectRepo.InsertProject(ctx, tx, project)

	if err != nil {
		panic(err)
	}

	projectMember := &models.ProjectMember{
		Id:        utils.GenerateUUID(),
		ProjectId: project.Id,
		UserId:    params.UserId,
		Role:      models.ProjectMemberRoleOwner,
		CreatedAt: time.Now().UTC(),
	}

	err = s.projectRepo.InsertProjectMember(ctx, tx, projectMember)

	if err != nil {
		panic(err)
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

	return project, nil
}

type GetProjectByIdentifierParams struct {
	Identifier string
	Type       models.ProjectType
}

func (s *projectService) GetProjectByIdentifier(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, params.Identifier, params.Type)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, errors.ErrProjectNotFound
	}

	return project, nil
}

type GetProjectMembersParams struct {
	Identifier string
	Type       models.ProjectType
}

func (s *projectService) GetProjectMembers(ctx context.Context, params GetProjectByIdentifierParams) ([]models.ProjectMember, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, params.Identifier, params.Type)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, errors.ErrProjectNotFound
	}

	members, err := s.projectRepo.FindProjectMembersByProjectId(ctx, s.db, project.Id)

	if err != nil {
		panic(err)
	}

	return members, nil
}
