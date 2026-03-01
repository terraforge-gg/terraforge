package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/handler"
	"github.com/terraforge-gg/terraforge/internal/lib/meilisearch"
	custom_middleware "github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

func NewServer(cfg *config.Config, logger *slog.Logger, db *sql.DB) (*echo.Echo, error) {

	jwtValidator, err := auth.NewValidator(cfg.AuthUrl + "/api/auth/jwks")
	authMiddleware := custom_middleware.JWTMiddleware(jwtValidator)
	authOptionalMiddleware := custom_middleware.OptionalJWTMiddleware(jwtValidator)

	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS validator: %w", err)
	}

	meiliClient := meilisearch.NewMeiliSearch(cfg.MeiliSearchHostUrl, cfg.MeiliSearchMasterKey)
	meiliSearchRepo := repository.NewMeiliSearchRepository(logger, meiliClient)
	searchService := service.NewSearchService(logger, meiliSearchRepo)

	authService := service.NewAuthService(logger, cfg.AuthUrl)

	projectRepo := repository.NewProjectRepository()
	projectService := service.NewProjectService(logger, db, projectRepo, meiliSearchRepo)
	projectHandler := handler.NewProjectHandler(cfg, logger, projectService, searchService)

	validate := validator.New()
	validate.RegisterValidation("url_slug", validation.ValidateUrlSlug)
	validate.RegisterValidation("project_type", validation.ValidateProjectType)

	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second),
		health.WithCheck(health.Check{
			Name:    "database",
			Timeout: 2 * time.Second,
			Check:   db.PingContext,
		}),
		health.WithCheck(health.Check{
			Name:    "search",
			Timeout: 2 * time.Second,
			Check:   meiliClient.Health,
		}),
		health.WithCheck(health.Check{
			Name:    "auth",
			Timeout: 2 * time.Second,
			Check:   authService.Health,
		}),
	)

	e := echo.New()

	e.Validator = &validation.Validator{Validator: validate}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendUrl},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	e.GET("/ready", echo.WrapHandler(health.NewHandler(checker)))

	v1 := e.Group("/v1")
	v1.File("/openapi.yml", "./docs/openapi.yml")
	v1.POST("/projects", projectHandler.CreateProject, authMiddleware)
	v1.GET("/projects", projectHandler.SearchProjects)
	v1.GET("/projects/:identifier", projectHandler.GetProjectByIdentifier, authOptionalMiddleware)
	v1.GET("/projects/:identifier/members", projectHandler.GetProjectMembers, authOptionalMiddleware)
	v1.PATCH("/projects/:identifier", projectHandler.UpdateProject, authMiddleware)
	v1.DELETE("/projects/:identifier", projectHandler.DeleteProject, authMiddleware)

	return e, nil
}
