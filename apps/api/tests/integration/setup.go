package integration

import (
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/cache"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/handler"
	"github.com/terraforge-gg/terraforge/internal/lib/aws"
	"github.com/terraforge-gg/terraforge/internal/logger"
	"github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

type testEnv struct {
	server *echo.Echo
	token1 string
	token2 string
	cfg    *config.Config
}

func newTestEnv(t *testing.T) *testEnv {
	db, err := database.NewTestDatabaseWithCleanup(t)
	testAuth, err := auth.NewTestAuth()
	require.NoError(t, err)

	jwtValidator := auth.NewValidatorFromJWKS(testAuth.JWKS)
	authMiddleware := middleware.JWTMiddleware(jwtValidator)
	authOptionalMiddleware := middleware.OptionalJWTMiddleware(jwtValidator)

	token1 := generateTestToken(t, testAuth, database.TestUser1Id, database.TestUser1Username, database.TestUser1Email)
	token2 := generateTestToken(t, testAuth, database.TestUser2Id, database.TestUser2Username, database.TestUser2Email)

	cfg := &config.Config{
		Env:         "test",
		DatabaseUrl: db.DatabaseUrl,
		R2Bucket:    "terraforge-app-test",
	}

	s3Client, host, err := aws.NewTestLocalStackS3Client(t, cfg.R2Bucket)
	cfg.CdnUrl = host

	log := logger.New()

	projectRepo := repository.NewProjectRepository()
	userRepo := repository.NewUserRepository()
	searchRepo := repository.NewMockSearchRepository()
	projectCache := cache.NewMockProjectCache()
	loaderVersionRepo := repository.NewLoaderVersionRepository()
	objectStoreService := service.NewObjectStoreService(s3Client, cfg.R2Bucket)

	projectService := service.NewProjectService(log, db.Db, projectRepo, searchRepo, projectCache, userRepo)
	searchService := service.NewMockSearchService()

	projectHandler := handler.NewProjectHandler(cfg, log, projectService, searchService)

	projectReleaseRepo := repository.NewProjectReleaseRepository()
	projectReleaseService := service.NewProjectReleaseService(
		log,
		cfg.CdnUrl,
		db.Db,
		projectRepo,
		projectReleaseRepo,
		loaderVersionRepo,
		objectStoreService,
	)
	projectReleaseHandler := handler.NewProjectReleaseHandler(cfg, log, projectReleaseService)

	validate := validation.NewValidator(cfg)

	e := echo.New()
	e.Validator = &validation.Validator{Validator: validate}

	v1 := e.Group("/v1")
	v1.POST("/projects", projectHandler.CreateProject, authMiddleware)
	v1.GET("/projects/:identifier", projectHandler.GetProjectByIdentifier, authOptionalMiddleware)
	v1.PATCH("/projects/:identifier", projectHandler.UpdateProject, authMiddleware)
	v1.DELETE("/projects/:identifier", projectHandler.DeleteProject, authMiddleware)
	v1.GET("/projects", projectHandler.SearchProjects)

	v1.POST("/projects/:identifier/releases", projectReleaseHandler.CreateRelease, authMiddleware)
	v1.GET("/projects/:identifier/releases", projectReleaseHandler.GetReleases, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/:releaseId", projectReleaseHandler.GetRelease, authOptionalMiddleware)
	v1.GET("/projects/:identifier/releases/upload-url", projectReleaseHandler.GeneratePresignedPutUrl, authMiddleware)

	return &testEnv{
		server: e,
		token1: token1,
		token2: token2,
		cfg:    cfg,
	}
}
