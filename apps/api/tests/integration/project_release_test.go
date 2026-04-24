package integration

import (
	"bytes"
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
	"github.com/terraforge-gg/terraforge/internal/utils"
)

func TestIntegration_CreateRelease(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)
	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)

	var uploadUrl string
	err := json.Unmarshal(uploadUrlRec.Body.Bytes(), &uploadUrl)
	require.NoError(t, err)

	fileContent := []byte("fake mod file content for testing")
	putReq, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(fileContent))
	require.NoError(t, err)
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putClient := &http.Client{}
	putRes, err := putClient.Do(putReq)
	require.NoError(t, err)
	defer putRes.Body.Close()
	assert.Equal(t, http.StatusOK, putRes.StatusCode)

	origin, pathname, err := utils.ExtractOriginAndPathFromUrl(uploadUrl)
	require.NoError(t, err)
	fileUrl := origin + pathname

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		fileUrl,
		nil,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)

	if releaseRec.Code != http.StatusOK {
		t.Logf("Response body: %s", releaseRec.Body.String())
	}
	assert.Equal(t, http.StatusOK, releaseRec.Code)

	var response dto.ProjectReleaseResponse
	err = json.Unmarshal(releaseRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, ExampleReleaseName, response.Name)
	assert.Equal(t, ExampleReleaseVersion, response.VersionNumber)
	assert.Equal(t, database.TestLoaderVersionId, response.LoaderVersion.Id)
}

func TestIntegration_CreateRelease_DuplicateVersionNumber(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)
	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)

	var uploadUrl string
	err := json.Unmarshal(uploadUrlRec.Body.Bytes(), &uploadUrl)
	require.NoError(t, err)

	fileContent := []byte("fake mod file content for testing")
	putReq, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(fileContent))
	require.NoError(t, err)
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putClient := &http.Client{}
	putRes, err := putClient.Do(putReq)
	require.NoError(t, err)
	defer putRes.Body.Close()
	assert.Equal(t, http.StatusOK, putRes.StatusCode)

	origin, pathname, err := utils.ExtractOriginAndPathFromUrl(uploadUrl)
	require.NoError(t, err)
	fileUrl := origin + pathname

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		fileUrl,
		nil,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)
	assert.Equal(t, http.StatusOK, releaseRec.Code)

	releaseReq2 := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq2.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec2 := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec2, releaseReq2)

	var response dto.ProblemDetails
	err = json.Unmarshal(releaseRec2.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response.Detail, "already exists")
	assert.Equal(t, http.StatusBadRequest, releaseRec2.Code)
}

func TestIntegration_GetReleases(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)
	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)

	var uploadUrl string
	err := json.Unmarshal(uploadUrlRec.Body.Bytes(), &uploadUrl)
	require.NoError(t, err)

	fileContent := []byte("fake mod file content for testing")
	putReq, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(fileContent))
	require.NoError(t, err)
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putClient := &http.Client{}
	putRes, err := putClient.Do(putReq)
	require.NoError(t, err)
	defer putRes.Body.Close()
	assert.Equal(t, http.StatusOK, putRes.StatusCode)

	origin, pathname, err := utils.ExtractOriginAndPathFromUrl(uploadUrl)
	require.NoError(t, err)
	fileUrl := origin + pathname

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		fileUrl,
		nil,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)
	assert.Equal(t, http.StatusOK, releaseRec.Code)

	getReleasesReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases",
		nil,
	)
	getReleasesReq.Header.Set("Authorization", "Bearer "+env.token1)
	getReleasesRec := httptest.NewRecorder()
	env.server.ServeHTTP(getReleasesRec, getReleasesReq)

	assert.Equal(t, http.StatusOK, getReleasesRec.Code)

	var response []dto.ProjectReleaseResponse
	err = json.Unmarshal(getReleasesRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, ExampleReleaseName, response[0].Name)
	assert.Equal(t, ExampleReleaseVersion, response[0].VersionNumber)
}

func TestIntegration_GetRelease(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)
	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)

	var uploadUrl string
	err := json.Unmarshal(uploadUrlRec.Body.Bytes(), &uploadUrl)
	require.NoError(t, err)

	fileContent := []byte("fake mod file content for testing")
	putReq, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(fileContent))
	require.NoError(t, err)
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putClient := &http.Client{}
	putRes, err := putClient.Do(putReq)
	require.NoError(t, err)
	defer putRes.Body.Close()
	assert.Equal(t, http.StatusOK, putRes.StatusCode)

	origin, pathname, err := utils.ExtractOriginAndPathFromUrl(uploadUrl)
	require.NoError(t, err)
	fileUrl := origin + pathname

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		fileUrl,
		nil,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)
	assert.Equal(t, http.StatusOK, releaseRec.Code)

	var createdRelease dto.ProjectReleaseResponse
	err = json.Unmarshal(releaseRec.Body.Bytes(), &createdRelease)
	require.NoError(t, err)

	getReleaseReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/"+createdRelease.Id,
		nil,
	)
	getReleaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	getReleaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(getReleaseRec, getReleaseReq)

	assert.Equal(t, http.StatusOK, getReleaseRec.Code)

	var response dto.ProjectReleaseResponse
	err = json.Unmarshal(getReleaseRec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, ExampleReleaseName, response.Name)
	assert.Equal(t, ExampleReleaseVersion, response.VersionNumber)
}

func TestIntegration_GeneratePresignedPutUrl(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)

	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)
	assert.NotEmpty(t, uploadUrlRec.Body.String())
}

func TestIntegration_CreateRelease_NonExistentDependency(t *testing.T) {
	env := newTestEnv(t)

	projectSummary := ExampleModSummary
	projectBody := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &projectSummary, "mod")

	projectReq := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(projectBody))
	projectReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	projectReq.Header.Set("Authorization", "Bearer "+env.token1)
	projectRec := httptest.NewRecorder()
	env.server.ServeHTTP(projectRec, projectReq)
	assert.Equal(t, http.StatusCreated, projectRec.Code)

	uploadUrlReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/projects/"+ExampleModSlug+"/releases/upload-url?fileSize="+ExampleReleaseFileSize,
		nil,
	)
	uploadUrlReq.Header.Set("Authorization", "Bearer "+env.token1)
	uploadUrlRec := httptest.NewRecorder()
	env.server.ServeHTTP(uploadUrlRec, uploadUrlReq)
	assert.Equal(t, http.StatusOK, uploadUrlRec.Code)

	var uploadUrl string
	err := json.Unmarshal(uploadUrlRec.Body.Bytes(), &uploadUrl)
	require.NoError(t, err)

	fileContent := []byte("fake mod file content for testing")
	putReq, err := http.NewRequest(http.MethodPut, uploadUrl, bytes.NewReader(fileContent))
	require.NoError(t, err)
	putReq.Header.Set("Content-Type", "application/octet-stream")
	putClient := &http.Client{}
	putRes, err := putClient.Do(putReq)
	require.NoError(t, err)
	defer putRes.Body.Close()
	assert.Equal(t, http.StatusOK, putRes.StatusCode)

	origin, pathname, err := utils.ExtractOriginAndPathFromUrl(uploadUrl)
	require.NoError(t, err)
	fileUrl := origin + pathname

	dependencies := []dto.CreateProjectReleaseRequestDependency{
		{
			ProjectId:        "nonexistent-project",
			MinVersionNumber: nil,
			Type:             "required",
		},
	}

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		fileUrl,
		dependencies,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseReq.Header.Set("Authorization", "Bearer "+env.token1)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)

	assert.Equal(t, http.StatusNotFound, releaseRec.Code)
}

func TestIntegration_CreateRelease_Unauthorized(t *testing.T) {
	env := newTestEnv(t)

	summary := ExampleModSummary
	body := createCreateProjectRequestBody(t, ExampleModName, ExampleModSlug, &summary, "mod")

	req := httptest.NewRequest(http.MethodPost, "/v1/projects", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+env.token1)
	rec := httptest.NewRecorder()
	env.server.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	changelog := ExampleReleaseChangelog
	releaseBody := createCreateReleaseRequestBody(
		t,
		ExampleReleaseName,
		ExampleReleaseVersion,
		&changelog,
		database.TestLoaderVersionId,
		env.cfg.CdnUrl+ExampleReleaseFilePath,
		nil,
	)

	releaseReq := httptest.NewRequest(
		http.MethodPost,
		"/v1/projects/"+ExampleModSlug+"/releases",
		strings.NewReader(releaseBody),
	)
	releaseReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	releaseRec := httptest.NewRecorder()
	env.server.ServeHTTP(releaseRec, releaseReq)

	assert.Equal(t, http.StatusUnauthorized, releaseRec.Code)
}
