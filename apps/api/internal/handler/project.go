package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/dto"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

type ProjectHandler struct {
	cfg            *config.Config
	logger         *slog.Logger
	projectService service.ProjectService
}

func NewProjectHandler(cfg *config.Config, logger *slog.Logger, projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		cfg:            cfg,
		logger:         logger,
		projectService: projectService,
	}
}

func (h *ProjectHandler) GetProjectByIdentifier(c *echo.Context) error {
	identifier := c.Param("identifier")
	_userId := c.Get("userId")
	ctx := c.Request().Context()

	var userId *string = nil
	if _userId != nil {
		s := _userId.(string)
		userId = &s
	}

	project, err := h.projectService.GetProjectByIdentifier(ctx, service.GetProjectByIdentifierParams{
		Identifier: identifier,
		UserId:     userId,
	})

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrProjectNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "Project not found.",
			})
		default:
			panic(err)
		}
	}

	return c.JSON(http.StatusOK, dto.ProjectToProjectResponse(*project))
}

func (h *ProjectHandler) GetProjectMembers(c *echo.Context) error {
	identifier := c.Param("identifier")
	_userId := c.Get("userId")
	ctx := c.Request().Context()

	var userId *string = nil
	if _userId != nil {
		s := _userId.(string)
		userId = &s
	}

	members, err := h.projectService.GetProjectMembers(ctx, service.GetProjectByIdentifierParams{
		Identifier: identifier,
		UserId:     userId,
	})

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrProjectNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "Project not found.",
			})
		default:
			panic(err)
		}
	}

	response := make([]dto.ProjectMemberResponse, len(members))

	for i, m := range members {
		response[i] = dto.ProjectMemberToProjectMemberResponse(m)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) CreateProject(c *echo.Context) error {
	userId := c.Get("userId").(string)
	var req dto.CreateProjectRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "Invalid request.",
		})
	}

	if err := c.Validate(&req); err != nil {
		if valErr, ok := err.(*validation.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "One or more fields failed validation.",
				Errors: valErr.Errors,
			})
		}

		panic(err)
	}

	ctx := c.Request().Context()
	project, err := h.projectService.CreateUserProject(ctx, service.CreateUserProjectParams{
		Name:    req.Name,
		Slug:    req.Slug,
		Summary: req.Summary,
		UserId:  userId,
		Type:    models.ProjectType(req.Type),
	})

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrProjectSlugUsed):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Project slug is not available.",
			})
		default:
			panic(err)
		}
	}

	return c.JSON(http.StatusCreated, dto.ProjectToProjectResponse(*project))
}
