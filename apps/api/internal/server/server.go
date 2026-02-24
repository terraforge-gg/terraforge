package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/handler"
	custom_middleware "github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

func NewServer(cfg *config.Config, logger *slog.Logger, db *sql.DB) (*echo.Echo, error) {
	e := echo.New()

	jwtValidator, err := auth.NewValidator(cfg.AuthUrl + "/api/auth/jwks")

	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS validator: %w", err)
	}

	projectRepo := repository.NewProjectRepository()
	projectService := service.NewProjectService(logger, db, projectRepo)
	projectHandler := handler.NewProjectHandler(cfg, logger, projectService)

	validate := validator.New()
	validate.RegisterValidation("url_slug", validation.ValidateUrlSlug)
	validate.RegisterValidation("project_type", validation.ValidateProjectType)

	e.Validator = &validation.Validator{Validator: validate}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"name": "terraforge",
			"env":  cfg.Env,
		})
	})

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"timestamp": time.Now().UTC().String(),
		})
	})

	v1 := e.Group("/v1")
	v1.POST("/projects", projectHandler.CreateProject, custom_middleware.JWTMiddleware(jwtValidator))

	return e, nil
}
