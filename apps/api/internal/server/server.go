package server

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/dto"
	"github.com/terraforge-gg/terraforge/internal/handler"
	"github.com/terraforge-gg/terraforge/internal/lib/aws"
	"github.com/terraforge-gg/terraforge/internal/lib/meilisearch"
	custom_middleware "github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/seed"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/utils"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

func NewServer(cfg *config.Config, logger *slog.Logger, db *sql.DB) (*echo.Echo, error) {
	jwtValidator, err := auth.NewValidator(cfg.AuthUrl + "/api/auth/jwks")
	authMiddleware := custom_middleware.JWTMiddleware(jwtValidator)
	authOptionalMiddleware := custom_middleware.OptionalJWTMiddleware(jwtValidator)

	if err != nil {
		return nil, fmt.Errorf("failed to create JWKS validator: %w", err)
	}

	aws_config, err := aws.NewAwsConfig(cfg)

	if err != nil {
		log.Fatal(err)
	}

	s3_client := aws.NewS3Client(cfg, aws_config)
	objectStoreService := service.NewObjectStoreService(s3_client, cfg.S3AssetsBucketName)

	meiliClient := meilisearch.NewMeiliSearch(cfg.MeiliSearchHostUrl, cfg.MeiliSearchMasterKey)
	meiliSearchRepo := repository.NewMeiliSearchRepository(logger, meiliClient)
	searchService := service.NewSearchService(logger, meiliSearchRepo)

	loaderVersionRepo := repository.NewLoaderVersionRepository()
	loaderVersionService := service.NewLoaderVersionService(logger, db, loaderVersionRepo)
	loaderVersionHandler := handler.NewLoaderVersionHandler(cfg, logger, loaderVersionService)

	projectRepo := repository.NewProjectRepository()
	projectService := service.NewProjectService(logger, db, projectRepo, meiliSearchRepo)
	projectHandler := handler.NewProjectHandler(cfg, logger, projectService, searchService)

	projectReleasenRepo := repository.NewProjectReleaseRepository()
	projectReleaseService := service.NewProjectReleaseService(logger, cfg.CdnUrl, db, projectRepo, projectReleasenRepo, loaderVersionRepo, objectStoreService)
	projectReleaseHandler := handler.NewProjectReleaseHandler(cfg, logger, projectReleaseService)

	if cfg.SeedDb {
		seed.SeedLoaderVersions(logger, loaderVersionService)
	}

	validate := validator.New()
	validate.RegisterValidation("url_slug", validation.ValidateUrlSlug)
	validate.RegisterValidation("project_type", validation.ValidateProjectType)
	validate.RegisterValidation("project_version_dependency_type", validation.ValidateProjectDependencyType)
	validate.RegisterValidation("file_url", validation.ValidateFileUrl)

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
		now := time.Now().UTC()
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"timestamp": now.Format(time.RFC3339),
		})
	})

	e.GET("/debug/protected", func(c *echo.Context) error {
		userId, ok := utils.GetUserId(c)

		if !ok {
			return c.JSON(http.StatusUnauthorized, dto.ProblemDetails{
				Title:  "Unauthorized",
				Status: http.StatusUnauthorized,
				Detail: "Unauthorized",
			})
		}

		return c.String(http.StatusOK, "Hello "+userId)
	}, authMiddleware)

	e.GET("/debug/protected-optional", func(c *echo.Context) error {
		userId, ok := utils.GetUserId(c)

		if !ok {
			userId = "NONE"
		}

		return c.String(http.StatusOK, "Hello "+userId)
	}, authOptionalMiddleware)

	v1 := e.Group("/v1")
	v1.File("/openapi.yml", "./docs/openapi.yml")

	v1.GET("/loader-versions/:id", loaderVersionHandler.GetLoaderVersionById)
	v1.GET("/loader-versions", loaderVersionHandler.GetLoaderVersions)

	v1.POST("/projects", projectHandler.CreateProject, authMiddleware)
	v1.GET("/projects/:identifier", projectHandler.GetProjectByIdentifier, authOptionalMiddleware)
	v1.GET("/projects/:identifier/members", projectHandler.GetProjectMembers, authOptionalMiddleware)

	v1.POST("/projects/:identifier/releases", projectReleaseHandler.CreateRelease, authMiddleware)
	v1.GET("/projects/:identifier/releases", projectReleaseHandler.GetReleases, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/:releaseId", projectReleaseHandler.GetRelease, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/upload-url", projectReleaseHandler.GeneratePresignedPutUrl, authMiddleware)

	return e, nil
}
