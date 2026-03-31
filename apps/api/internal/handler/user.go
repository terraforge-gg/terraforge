package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/dto"
	custom_errors "github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

type UserHandler struct {
	cfg            *config.Config
	logger         *slog.Logger
	projectService service.ProjectService
}

func NewUserHandler(cfg *config.Config, logger *slog.Logger, projectService service.ProjectService) *UserHandler {
	return &UserHandler{
		cfg:            cfg,
		logger:         logger,
		projectService: projectService,
	}
}

func (h *UserHandler) GetProjectsByUserId(c *echo.Context) error {
	ctx := c.Request().Context()
	userIdentifier := c.Param("userIdentifier")
	userId, _ := utils.GetSessionUserId(c)

	projects, err := h.projectService.GetProjectsByUserIdentifier(ctx, service.GetProjectsByUserIdentifierParams{
		UserIdentifier: userIdentifier,
		SessionUserId:  userId,
	})

	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "User not found.",
			})
		case errors.Is(err, custom_errors.ErrNoProjectsFound):
			return c.JSON(http.StatusNotFound, dto.ProblemDetails{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "No projects found.",
			})
		default:
			h.logger.Error("Unhandled get user projects error", "User identifier: ", userIdentifier, "Error:", err)
			return c.JSON(http.StatusInternalServerError, dto.ProblemDetails{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
			})
		}
	}

	userProjects := make([]dto.ProjectResponse, len(projects))

	for i, p := range projects {
		userProjects[i] = dto.ProjectToProjectResponse(p)
	}

	return c.JSON(http.StatusOK, userProjects)
}
