package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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
	searchService  service.SearchService
}

func NewProjectHandler(cfg *config.Config, logger *slog.Logger, projectService service.ProjectService, searchService service.SearchService) *ProjectHandler {
	return &ProjectHandler{
		cfg:            cfg,
		logger:         logger,
		projectService: projectService,
		searchService:  searchService,
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

func (h *ProjectHandler) UpdateProject(c *echo.Context) error {
	identifier := c.Param("identifier")
	userId := c.Get("userId").(string)
	ctx := c.Request().Context()

	var req dto.UpdateProjectRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
			Title:  "Bad Request",
			Status: http.StatusBadRequest,
			Detail: "Invalid request.",
		})
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	project, err := h.projectService.UpdateProject(ctx, service.UpdateProjectParams{
		Identifier:  identifier,
		Slug:        req.Slug,
		Name:        req.Name,
		Summary:     req.Summary,
		IconUrl:     req.IconUrl,
		Description: req.Description,
		UserId:      userId,
	})

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrProjectNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "Project not found.",
			})
		case errors.Is(err, custom_errors.ErrProjectUnauthorisedAction):
			return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
				Title:  "Unauthorised",
				Status: http.StatusUnauthorized,
				Detail: "You are not authorised to perform this action",
			})
		default:
			panic(err)
		}
	}

	return c.JSON(http.StatusOK, dto.ProjectToProjectResponse(*project))
}

func (h *ProjectHandler) DeleteProject(c *echo.Context) error {
	identifier := c.Param("identifier")
	userId := c.Get("userId").(string)
	ctx := c.Request().Context()

	err := h.projectService.DeleteProject(ctx, service.DeleteProjectParams{
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
		case errors.Is(err, custom_errors.ErrProjectUnauthorisedAction):
			return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
				Title:  "Unauthorised",
				Status: http.StatusUnauthorized,
				Detail: "You are not authorised to perform this action",
			})
		default:
			panic(err)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *ProjectHandler) SearchProjects(c *echo.Context) error {
	query := c.QueryParam("query")
	projectType := c.QueryParamOr("type", string(models.ProjectTypeMod))
	ctx := c.Request().Context()

	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err := strconv.ParseInt(c.QueryParam("offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	const maxLimit int64 = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	projects, totalHits := h.searchService.SearchProjects(ctx, query, projectType, limit, offset)

	response := dto.ProjectToProjectSearchResponse(projects, totalHits, limit, offset)

	return c.JSON(http.StatusOK, response)
}
