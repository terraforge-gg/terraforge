package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terraforge-gg/terraforge/internal/auth"
	"github.com/terraforge-gg/terraforge/internal/dto"
)

func generateTestToken(t *testing.T, testAuth *auth.TestAuth, userId string, username string, email string) string {
	t.Helper()
	token, err := testAuth.GenerateToken(userId, username, email)
	require.NoError(t, err)
	return token
}

const (
	ExampleModName    = "Example Mod"
	ExampleModSlug    = "example-mod"
	ExampleModSummary = "This is an example mod"

	CoolModName        = "Cool Mod"
	CoolModSlug        = "cool-mod"
	CoolModSummary     = "This is a cool mod"
	CoolModDescription = "This is a very cool mod that does some stuff"
	CoolModIconUrl     = "https://cdn.terraforge.gg/icons/cool-mod.png"

	DependencyModName    = "Dependency Mod"
	DependencyModSlug    = "dependency-mod"
	DependencyModSummary = "This is a dependency mod"
)

func createCreateProjectRequestBody(t *testing.T, name string, slug string, summary *string, projectType string) string {
	t.Helper()
	req := dto.CreateProjectRequest{
		Name:    name,
		Slug:    slug,
		Summary: summary,
		Type:    projectType,
	}

	b, err := json.Marshal(req)

	if err != nil {
		panic(err)
	}

	return string(b)
}

func createUpdateProjectRequestBody(t *testing.T, name *string, slug *string, summary *string, description *string, iconUrl *string) string {
	t.Helper()
	req := dto.UpdateProjectRequest{
		Name:        name,
		Slug:        slug,
		Summary:     summary,
		Description: description,
		IconUrl:     iconUrl,
	}

	b, err := json.Marshal(req)

	if err != nil {
		panic(err)
	}

	return string(b)
}

const (
	ExampleReleaseName      = "Example Release"
	ExampleReleaseVersion   = "1.0.0"
	ExampleReleaseChangelog = "Initial release"
	ExampleReleaseFilePath  = "/uploads/temp/1/test_123.tmod"
	ExampleReleaseFileSize  = "1024"
	CoolReleaseName         = "Cool Release"
	CoolReleaseVersion      = "1.1.0"
	CoolReleaseChangelog    = "Bug fixes and improvements"
)

func createCreateReleaseRequestBody(t *testing.T, name string, versionNumber string, changelog *string, loaderVersionId string, fileUrl string, dependencies []dto.CreateProjectReleaseRequestDependency) string {
	t.Helper()
	req := dto.CreateProjectReleaseRequest{
		Name:            name,
		VersionNumber:   versionNumber,
		Changelog:       changelog,
		LoaderVersionId: loaderVersionId,
		FileUrl:         fileUrl,
		Dependencies:    dependencies,
	}

	b, err := json.Marshal(req)

	if err != nil {
		panic(err)
	}

	return string(b)
}
