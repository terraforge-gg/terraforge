package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/cache"
	"github.com/terraforge-gg/terraforge/internal/config"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/dto"
	"github.com/terraforge-gg/terraforge/internal/handler"
	"github.com/terraforge-gg/terraforge/internal/logger"
	"github.com/terraforge-gg/terraforge/internal/middleware"
	"github.com/terraforge-gg/terraforge/internal/repository"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

var testServer *echo.Echo
var token1 string

const createExampleModRequestBody = `{
		"name": "Example Mod",
		"slug": "example-mod",
		"summary": "This is an example mod",
		"type": "mod"
	}`
const exampleModSlug = "example-mod"
const nonexistentProjectSlug = "non-existent"

func TestMain(m *testing.M) {
	t := &testing.T{}
	db, err := database.NewTestDatabaseWithCleanup(t)
	testAuth, err := auth.NewTestAuth()
	require.NoError(t, err)

	jwtValidator := auth.NewValidatorFromJWKS(testAuth.JWKS)
	authMiddleware := middleware.JWTMiddleware(jwtValidator)
	authOptionalMiddleware := middleware.OptionalJWTMiddleware(jwtValidator)

	token1 = generateTestToken(t, testAuth, database.TestUser1Id, database.TestUser1Username, database.TestUser1Email)

	cfg := &config.Config{
		Env:         "test",
		DatabaseUrl: db.DatabaseUrl,
	}

	logger := logger.New()

	projectRepo := repository.NewProjectRepository()
	userRepo := repository.NewUserRepository()
	searchRepo := repository.NewMockSearchRepository()
	projectCache := cache.NewMockProjectCache()

	projectService := service.NewProjectService(logger, db.Db, projectRepo, searchRepo, projectCache, userRepo)
	searchService := service.NewMockSearchService()

	projectHandler := handler.NewProjectHandler(cfg, logger, projectService, searchService)

	validate := validation.NewValidator()

	e := echo.New()
	e.Validator = &validation.Validator{Validator: validate}

	v1 := e.Group("/v1")
	v1.POST("/projects", projectHandler.CreateProject, authMiddleware)
	v1.GET("/projects/:identifier", projectHandler.GetProjectByIdentifier, authOptionalMiddleware)
	v1.PATCH("/projects/:identifier", projectHandler.UpdateProject, authMiddleware)
	v1.DELETE("/projects/:identifier", projectHandler.DeleteProject, authMiddleware)
	v1.GET("/projects", projectHandler.SearchProjects)

	testServer = e
	os.Exit(m.Run())
}

func TestIntegration_CreateProject(t *testing.T) {
	body := createExampleModRequestBody

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()

	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response dto.ProjectResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "Example Mod", response.Name)
	assert.Equal(t, "example-mod", response.Slug)
	assert.Equal(t, "This is an example mod", *response.Summary)
	assert.Equal(t, database.TestUser1Id, response.UserId)
	assert.Equal(t, "mod", response.Type)
	assert.Equal(t, "draft", response.Status)
}

func TestIntegration_CreateProject_Validation(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{"EmptyName", `{"slug":"example-mod-1","type":"mod"}`, http.StatusBadRequest},
		{"EmptySlug", `{"name":"Example Mod 2","type":"mod"}`, http.StatusBadRequest},
		{"EmptyType", `{"name":"Example Mod 3","slug":"example-mod-3"}`, http.StatusBadRequest},
		{"InvalidType", `{"name":"Example Mod 4","slug":"example-mo-4","type":"invalid"}`, http.StatusBadRequest},
		{"NameTooShort", `{"name":"E","slug":"example-mod-5","type":"mod"}`, http.StatusBadRequest},
		{"NameTooLong", `{"name":"Example Mod Example Mod Example Mod Example Mod 
			Example Mod Example Mod Example Mod Example Mod Example Mod Example Mod Example Mod Example Mod Example Mod",
			"slug":"example-mod-5","type":"mod"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("Authorization", "Bearer "+token1)
			rec := httptest.NewRecorder()
			testServer.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestIntegration_CreateProject_DuplicateSlug(t *testing.T) {
	body := createExampleModRequestBody

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIntegration_CreateProject_Unauthorized(t *testing.T) {
	body := createExampleModRequestBody

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestIntegration_GetProjectBySlug(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+exampleModSlug, nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	getRec := httptest.NewRecorder()
	testServer.ServeHTTP(getRec, req)

	assert.Equal(t, http.StatusOK, getRec.Code)

	var response dto.ProjectResponse
	err := json.Unmarshal(getRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Example Mod", response.Name)
	assert.Equal(t, "example-mod", response.Slug)
	assert.Equal(t, "This is an example mod", *response.Summary)
	assert.Equal(t, database.TestUser1Id, response.UserId)
	assert.Equal(t, "mod", response.Type)
	assert.Equal(t, "draft", response.Status)
}

func TestIntegration_GetProjectBySlug_Unauthenticated(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+exampleModSlug, nil)
	getRec := httptest.NewRecorder()
	testServer.ServeHTTP(getRec, req)

	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestIntegration_UpdateProject(t *testing.T) {
	updateBody := `{"name": "Updated Name", "slug": "updated-mod", "summary": "Updated summary"}`
	updateReq := httptest.NewRequest(http.MethodPatch, "/v1/projects/"+exampleModSlug, strings.NewReader(updateBody))
	updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	updateReq.Header.Set("Authorization", "Bearer "+token1)
	updateRec := httptest.NewRecorder()
	testServer.ServeHTTP(updateRec, updateReq)

	assert.Equal(t, http.StatusOK, updateRec.Code)

	var response dto.ProjectResponse
	err := json.Unmarshal(updateRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Name)
	assert.Equal(t, "updated-mod", response.Slug)

	getOldReq := httptest.NewRequest(http.MethodGet, "/v1/projects/"+exampleModSlug, nil)
	getOldRec := httptest.NewRecorder()
	testServer.ServeHTTP(getOldRec, getOldReq)

	assert.Equal(t, http.StatusNotFound, getOldRec.Code)
}

func TestIntegration_DeleteProject(t *testing.T) {
	body := createExampleModRequestBody

	createReq := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createReq.Header.Set("Authorization", "Bearer "+token1)
	createRec := httptest.NewRecorder()
	testServer.ServeHTTP(createRec, createReq)
	require.Equal(t, http.StatusCreated, createRec.Code)

	deleteReq := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+exampleModSlug, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token1)
	deleteRec := httptest.NewRecorder()
	testServer.ServeHTTP(deleteRec, deleteReq)

	assert.Equal(t, http.StatusOK, deleteRec.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/projects/"+exampleModSlug, nil)
	getRec := httptest.NewRecorder()
	testServer.ServeHTTP(getRec, getReq)

	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestIntegration_GetProject_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+nonexistentProjectSlug, nil)
	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIntegration_UpdateProject_NotFound(t *testing.T) {
	body := `{"name": "Updated"}`
	req := httptest.NewRequest(http.MethodPatch, "/v1/projects/"+nonexistentProjectSlug, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestIntegration_DeleteProject_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+nonexistentProjectSlug, nil)
	req.Header.Set("Authorization", "Bearer "+token1)
	rec := httptest.NewRecorder()
	testServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
