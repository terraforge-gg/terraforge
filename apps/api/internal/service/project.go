package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/terraforge-gg/terraforge/internal/database"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

type ProjectService interface {
	CreateUserProject(ctx context.Context, params CreateUserProjectParams) (*models.Project, error)
	GetProjectByIdentifier(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error)
	GetProjectMembers(ctx context.Context, params GetProjectByIdentifierParams) ([]models.ProjectMember, error)
	UpdateProject(ctx context.Context, params UpdateProjectParams) (*models.Project, error)
	DeleteProject(ctx context.Context, params DeleteProjectParams) error
}

type projectService struct {
	logger      *slog.Logger
	db          *sql.DB
	projectRepo repository.ProjectRepository
	searchRepo  repository.SearchRepository
}

func NewProjectService(logger *slog.Logger, db *sql.DB, projectRepo repository.ProjectRepository, searchRepo repository.SearchRepository) ProjectService {
	return &projectService{logger: logger, db: db, projectRepo: projectRepo, searchRepo: searchRepo}
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

	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err != nil {
		panic(err)
	}

	err = s.projectRepo.InsertProject(ctx, tx, project)

	if err != nil {
		switch {
		case errors.Is(err, database.ErrUniqueViolation):
			return nil, custom_errors.ErrProjectSlugUsed
		}
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

	go func() {
		err = s.searchRepo.IndexProject(context.Background(), project)
		if err != nil {
			s.logger.Error("Failed to index project.", "Project: ", project, "Error: ", err)
		}
	}()

	return project, nil
}

type GetProjectByIdentifierParams struct {
	Identifier string
	UserId     *string
}

func (s *projectService) GetProjectByIdentifier(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error) {
	project, err := s.projectRepo.FindProjectByIdentifier(ctx, s.db, params.Identifier, params.UserId)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	return project, nil
}

type GetProjectMembersParams struct {
	Identifier string
	UserId     *string
}

func (s *projectService) GetProjectMembers(ctx context.Context, params GetProjectByIdentifierParams) ([]models.ProjectMember, error) {
	members, err := s.projectRepo.FindProjectMembersByProjectIdentifier(ctx, s.db, params.Identifier, params.UserId)

	if err != nil {
		panic(err)
	}

	if members == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	return members, nil
}

type UpdateProjectParams struct {
	Identifier  string
	Name        *string
	Slug        *string
	Summary     *string
	Description *string
	IconUrl     *string
	UserId      string
}

func (s *projectService) UpdateProject(ctx context.Context, params UpdateProjectParams) (*models.Project, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, params.Identifier, &params.UserId)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, s.db, project.Id, params.UserId)

	if err != nil {
		panic(err)
	}

	if projectMember == nil {
		return nil, custom_errors.ErrProjectNotFound
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return nil, custom_errors.ErrProjectUnauthorisedAction
	}

	if params.Name != nil {
		project.Name = *params.Name
	}

	if params.Slug != nil {
		project.Slug = *params.Slug
	}

	if params.Summary != nil {
		if *params.Summary == "" {
			project.Summary = nil
		} else {
			project.Summary = params.Summary
		}
	}

	if params.Description != nil {
		if *params.Description == "" {
			project.Description = nil
		} else {
			project.Description = params.Description
		}
	}

	if params.IconUrl != nil {
		if *params.IconUrl == "" {
			project.IconUrl = nil
		} else {
			project.IconUrl = params.IconUrl
		}
	}

	err = s.projectRepo.UpdateProject(ctx, tx, *project)

	if err != nil {
		panic(err)
	}

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

	go func() {
		err = s.searchRepo.UpdateProject(context.Background(), project)
		if err != nil {
			s.logger.Error("Failed to update project document.", "Project: ", project, "Error: ", err)
		}
	}()

	return project, nil
}

type DeleteProjectParams struct {
	Identifier string
	UserId     string
}

func (s *projectService) DeleteProject(ctx context.Context, params DeleteProjectParams) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	project, err := s.projectRepo.FindProjectByIdentifier(ctx, tx, params.Identifier, &params.UserId)

	if err != nil {
		panic(err)
	}

	if project == nil {
		return custom_errors.ErrProjectNotFound
	}

	projectMember, err := s.projectRepo.FindProjectMemberByProjectIdAndUserId(ctx, s.db, project.Id, params.UserId)

	if err != nil {
		panic(err)
	}

	if projectMember == nil {
		return custom_errors.ErrProjectNotFound
	}

	if projectMember.Role != models.ProjectMemberRoleOwner {
		return custom_errors.ErrProjectUnauthorisedAction
	}

	deletedAt := time.Now().UTC()
	err = s.projectRepo.DeleteProjectByIdentifier(ctx, tx, project.Id, deletedAt)

	err = tx.Commit()

	if err != nil {
		panic(err)
	}

	go func() {
		err = s.searchRepo.DeleteProject(context.Background(), project.Id)
		if err != nil {
			s.logger.Error("Failed to delete project document.", "Project: ", project, "Error: ", err)
		}
	}()

	return nil
}
