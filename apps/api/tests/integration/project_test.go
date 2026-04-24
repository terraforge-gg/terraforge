package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terraforge-gg/terraforge/internal/database"
	"github.com/terraforge-gg/terraforge/internal/dto"
	"github.com/terraforge-gg/terraforge/internal/models"
)

func TestIntegration_CreateProject(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))

	// Act
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	var response dto.ProjectResponse
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, ExampleModName, response.Name)
	assert.Equal(t, ExampleModSlug, response.Slug)
	assert.Equal(t, summary, *response.Summary)
	assert.Equal(t, database.TestUser1Id, response.UserId)
	assert.Equal(t, string(string(models.ProjectTypeMod)), response.Type)
	assert.Equal(t, string(models.ProjectStatusDraft), response.Status)
}

func TestIntegration_CreateProject_DuplicateSlug(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Act
	env.server.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)

}

func TestIntegration_CreateProject_Unauthorized(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))

	// Act
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestIntegration_GetDraftProjectBySlug_Authenticated(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Act
	req = httptest.NewRequest(http.MethodGet, "/v1/projects/"+ExampleModSlug, nil)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	getRec := httptest.NewRecorder()
	env.server.ServeHTTP(getRec, req)

	assert.Equal(t, http.StatusOK, getRec.Code)

	var response dto.ProjectResponse
	err := json.Unmarshal(getRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, ExampleModName, response.Name)
	assert.Equal(t, ExampleModSlug, response.Slug)
	assert.Equal(t, summary, *response.Summary)
	assert.Equal(t, database.TestUser1Id, response.UserId)
	assert.Equal(t, string(string(models.ProjectTypeMod)), response.Type)
	assert.Equal(t, string(models.ProjectStatusDraft), response.Status)
}

func TestIntegration_GetDraftProjectBySlug_Unauthenticated(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Act
	req = httptest.NewRequest(http.MethodGet, "/v1/projects/"+ExampleModSlug, nil)
	req.Header.Set("Authorization", "")
	getRec := httptest.NewRecorder()
	env.server.ServeHTTP(getRec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestIntegration_UpdateProject(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Act
	newName := CoolModName
	newSlug := CoolModSlug
	newSummary := CoolModSummary
	newDescription := CoolModDescription
	newIconUrl := CoolModIconUrl

	updateBody := createUpdateProjectRequestBody(t, &newName, &newSlug, &newSummary, &newDescription, &newIconUrl)
	updateReq := httptest.NewRequest(http.MethodPatch, "/v1/projects/"+ExampleModSlug, strings.NewReader(updateBody))
	updateReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	updateReq.Header.Set("Authorization", "Bearer "+env.token1)
	updateRec := httptest.NewRecorder()
	env.server.ServeHTTP(updateRec, updateReq)

	// Assert
	assert.Equal(t, http.StatusOK, updateRec.Code)

	var response dto.ProjectResponse
	err := json.Unmarshal(updateRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, newName, response.Name)
	assert.Equal(t, newSlug, response.Slug)
	assert.Equal(t, newSummary, *response.Summary)
	assert.Equal(t, newDescription, *response.Description)
	assert.Equal(t, newIconUrl, *response.IconUrl)

	getOldReq := httptest.NewRequest(http.MethodGet, "/v1/projects/"+ExampleModSlug, nil)
	getOldRec := httptest.NewRecorder()
	env.server.ServeHTTP(getOldRec, getOldReq)

	assert.Equal(t, http.StatusNotFound, getOldRec.Code)
}

func TestIntegration_DeleteProject(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, string(models.ProjectTypeMod))
	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Act
	deleteReq := httptest.NewRequest(http.MethodDelete, "/v1/projects/"+ExampleModSlug, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+env.token1)
	deleteRec := httptest.NewRecorder()
	env.server.ServeHTTP(deleteRec, deleteReq)

	assert.Equal(t, http.StatusOK, deleteRec.Code)

	getReq := httptest.NewRequest(http.MethodGet, "/v1/projects/"+ExampleModSlug, nil)
	getRec := httptest.NewRecorder()
	env.server.ServeHTTP(getRec, getReq)

	assert.Equal(t, http.StatusNotFound, getRec.Code)
}

func TestIntegration_GetProject_NotFound(t *testing.T) {
	// Arrange
	env := newTestEnv(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/projects/"+CoolModSlug, nil)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()

	// Act
	env.server.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
