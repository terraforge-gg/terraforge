package server

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/cache"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/handler"
	"github.com/terraforge-gg/terraforge/internal/lib/aws"
	"github.com/terraforge-gg/terraforge/internal/lib/meilisearch"
	"github.com/terraforge-gg/terraforge/internal/lib/redis"
	custom_middleware "github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/seed"
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

	aws_config, err := aws.NewAwsConfig(cfg)

	if err != nil {
		log.Fatal(err)
	}

	s3_client := aws.NewS3Client(cfg, aws_config)
	objectStoreService := service.NewObjectStoreService(s3_client, cfg.R2Bucket)

	redisClient, err := redis.NewRedisClient(cfg.RedisUrl, cfg.RedisPassword)

	if err != nil {
		return nil, fmt.Errorf("failed to create redis client: %w", err)
	}

	projectCache := cache.NewProjectCache(redisClient)

	meiliClient := meilisearch.NewMeiliSearch(cfg.MeiliSearchHostUrl, cfg.MeiliSearchMasterKey)
	meiliSearchRepo := repository.NewMeiliSearchRepository(logger, meiliClient)
	searchService := service.NewSearchService(logger, meiliSearchRepo)

	authHealthCheckService := auth.NewAuthHealthCheckService(logger, cfg.AuthUrl)

	generalLimiter := custom_middleware.RateLimiter(redisClient.Client, custom_middleware.RateLimitGeneral)
	writeLimiter := custom_middleware.RateLimiter(redisClient.Client, custom_middleware.RateLimitWrite)
	searchLimiter := custom_middleware.RateLimiter(redisClient.Client, custom_middleware.RateLimitSearch)

	loaderVersionRepo := repository.NewLoaderVersionRepository()
	loaderVersionService := service.NewLoaderVersionService(logger, db, loaderVersionRepo)
	loaderVersionHandler := handler.NewLoaderVersionHandler(cfg, logger, loaderVersionService)

	userRepository := repository.NewUserRepository()

	projectRepo := repository.NewProjectRepository()
	projectService := service.NewProjectService(logger, db, projectRepo, meiliSearchRepo, projectCache, userRepository)
	projectHandler := handler.NewProjectHandler(cfg, logger, projectService, searchService)

	userHandler := handler.NewUserHandler(cfg, logger, projectService)

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
	validate.RegisterValidation("semver", validation.ValidateSemVer)

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
			Check:   authHealthCheckService.Health,
		}),
		health.WithCheck(health.Check{
			Name:    "redis",
			Timeout: 2 * time.Second,
			Check:   redisClient.Health,
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
		now := time.Now().UTC()
		return c.JSON(http.StatusOK, map[string]string{
			"status":    "ok",
			"timestamp": now.Format(time.RFC3339),
		})
	})

	e.GET("/ready", echo.WrapHandler(health.NewHandler(checker)))

	v1 := e.Group("/v1")
	v1.Use(generalLimiter)
	v1.File("/openapi.yml", "./docs/openapi.yml")

	v1.GET("/loader-versions/:id", loaderVersionHandler.GetLoaderVersionById)
	v1.GET("/loader-versions", loaderVersionHandler.GetLoaderVersions)

	v1.GET("/users/:userIdentifier/projects", userHandler.GetProjectsByUserId, authOptionalMiddleware)

	v1.POST("/projects", projectHandler.CreateProject, authMiddleware, writeLimiter)
	v1.GET("/projects", projectHandler.SearchProjects, searchLimiter)
	v1.GET("/projects/:identifier", projectHandler.GetProjectByIdentifier, authOptionalMiddleware)
	v1.GET("/projects/:identifier/members", projectHandler.GetProjectMembers, authOptionalMiddleware)
	v1.PATCH("/projects/:identifier", projectHandler.UpdateProject, authMiddleware, writeLimiter)
	v1.DELETE("/projects/:identifier", projectHandler.DeleteProject, authMiddleware, writeLimiter)

	v1.POST("/projects/:identifier/releases", projectReleaseHandler.CreateRelease, authMiddleware, writeLimiter)
	v1.GET("/projects/:identifier/releases", projectReleaseHandler.GetReleases, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/:releaseId", projectReleaseHandler.GetRelease, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/upload-url", projectReleaseHandler.GeneratePresignedPutUrl, authMiddleware)

	return e, nil
}
