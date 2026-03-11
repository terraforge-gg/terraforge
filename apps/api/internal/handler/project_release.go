package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/dto"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

type ProjectReleaseHandler struct {
	cfg                   *config.Config
	logger                *slog.Logger
	projectReleaseService service.ProjectReleaseService
}

func NewHandler(cfg *config.Config, logger *slog.Logger, projectReleaseService service.ProjectReleaseService) *ProjectReleaseHandler {
	return &ProjectReleaseHandler{
		cfg:                   cfg,
		logger:                logger,
		projectReleaseService: projectReleaseService,
	}
}

func (h *ProjectReleaseHandler) CreateRelease(c *echo.Context) error {
	identifier := c.Param("identifier")
	session := c.Get("session").(*auth.Session)
	ctx := c.Request().Context()

	var req dto.CreateProjectReleaseRequest

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

	deps := make([]service.CreateProjectReleaseDependencyParams, len(req.Dependencies))

	for i, dep := range req.Dependencies {
		deps[i] = service.CreateProjectReleaseDependencyParams{
			ProjectId:        dep.ProjectId,
			MinVersionNumber: dep.MinVersionNumber,
			Type:             dep.Type,
		}
	}

	version, err := h.projectReleaseService.CreateRelease(ctx, identifier, session.User.Id, service.CreateReleaseParams{
		Name:            req.Name,
		VersionNumber:   req.VersionNumber,
		Changelog:       req.Changelog,
		LoaderVersionId: req.LoaderVersionId,
		FileUrl:         req.FileUrl,
		Dependencies:    deps,
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
		case errors.Is(err, custom_errors.ErrLoaderVersionNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusBadRequest,
				Detail: "loader version '" + req.LoaderVersionId + "' not found.",
			})
		case errors.Is(err, custom_errors.ErrProjectReleaseDependencyNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusBadRequest,
				Detail: "Dependency project not found.",
			})
		case errors.Is(err, custom_errors.ErrCircularProjectReleaseDependency):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "A project cannot have a dependency on itself.",
			})
		case errors.Is(err, custom_errors.ErrDuplicateProjectReleaseDependency):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Duplicate dependencies cannot be added.",
			})
		case errors.Is(err, custom_errors.ErrProjectReleaseNumberAlreadyExists):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Project version" + " '" + req.VersionNumber + "' " + "already exists.",
			})
		case errors.Is(err, custom_errors.ErrProjectReleaseDependencyMinVersionDoesNotExist):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "A dependency project version with a supplied min version was not found.",
			})
		default:
			panic(err)
		}
	}
	versionDto := dto.MapToProjectReleaseResponse(*version)

	return c.JSON(http.StatusOK, versionDto)
}

func (h *ProjectReleaseHandler) GetReleases(c *echo.Context) error {
	identifier := c.Param("identifier")
	ctx := c.Request().Context()

	versions, err := h.projectReleaseService.GetReleasesByProjectId(ctx, identifier)

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

	versionsResponse := make([]dto.ProjectReleaseResponse, len(versions))

	for i, v := range versions {
		versionsResponse[i] = dto.MapToProjectReleaseResponse(v)
	}

	return c.JSON(http.StatusOK, versionsResponse)
}

func (h *ProjectReleaseHandler) GeneratePresignedPutUrl(c *echo.Context) error {
	identifier := c.Param("identifier")
	session := c.Get("session").(*auth.Session)
	fileSize := c.QueryParam("fileSize")
	ctx := c.Request().Context()

	url, err := h.projectReleaseService.GenerateProjectReleasePresignedPutUrl(ctx, identifier, session.User.Id, fileSize)

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
		case errors.Is(err, custom_errors.ErrProjectReleaseInvalidFileSize):
			return c.JSON(http.StatusBadRequest, dto.ProblemDetails{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "Invalid project version file size",
			})
		default:
			panic(err)
		}
	}

	return c.JSON(http.StatusOK, url)
}
