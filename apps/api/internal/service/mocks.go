package service

import (
	"context"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type MockProjectService struct {
	CreateUserProjectFunc           func(ctx context.Context, params CreateUserProjectParams) (*models.Project, error)
	GetProjectByIdentifierFunc      func(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error)
	GetProjectMembersFunc           func(ctx context.Context, params GetProjectMembersParams) ([]models.ProjectMember, error)
	UpdateProjectFunc               func(ctx context.Context, params UpdateProjectParams) (*models.Project, error)
	DeleteProjectFunc               func(ctx context.Context, params DeleteProjectParams) error
	GetProjectsByUserIdentifierFunc func(ctx context.Context, params GetProjectsByUserIdentifierParams) ([]models.Project, error)
}

func NewMockProjectService() *MockProjectService {
	return &MockProjectService{}
}

func (m *MockProjectService) CreateUserProject(ctx context.Context, params CreateUserProjectParams) (*models.Project, error) {
	if m.CreateUserProjectFunc != nil {
		return m.CreateUserProjectFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) GetProjectByIdentifier(ctx context.Context, params GetProjectByIdentifierParams) (*models.Project, error) {
	if m.GetProjectByIdentifierFunc != nil {
		return m.GetProjectByIdentifierFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) GetProjectMembers(ctx context.Context, params GetProjectMembersParams) ([]models.ProjectMember, error) {
	if m.GetProjectMembersFunc != nil {
		return m.GetProjectMembersFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) UpdateProject(ctx context.Context, params UpdateProjectParams) (*models.Project, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) DeleteProject(ctx context.Context, params DeleteProjectParams) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, params)
	}
	return nil
}

func (m *MockProjectService) GetProjectsByUserIdentifier(ctx context.Context, params GetProjectsByUserIdentifierParams) ([]models.Project, error) {
	if m.GetProjectsByUserIdentifierFunc != nil {
		return m.GetProjectsByUserIdentifierFunc(ctx, params)
	}
	return nil, nil
}

type MockSearchService struct {
	SearchProjectsFunc func(ctx context.Context, query string, projectType string, limit int64, offset int64) ([]models.Project, int64, error)
	HealthFunc         func(ctx context.Context) error
}

func NewMockSearchService() *MockSearchService {
	return &MockSearchService{}
}

func (m *MockSearchService) SearchProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) ([]models.Project, int64, error) {
	if m.SearchProjectsFunc != nil {
		return m.SearchProjectsFunc(ctx, query, projectType, limit, offset)
	}
	return nil, 0, nil
}

func (m *MockSearchService) Health(ctx context.Context) error {
	if m.HealthFunc != nil {
		return m.HealthFunc(ctx)
	}
	return nil
}

type MockProjectReleaseService struct {
	GetReleaseByIdWithDependenciesFunc func(ctx context.Context, projectIdentifier string, releaseId string, userId string) (*models.ProjectRelease, error)
	CreateReleaseFunc                  func(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error)
	GeneratePresignedPutUrlFunc        func(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error)
	GetReleasesByProjectIdFunc         func(ctx context.Context, id string, userId string) ([]models.ProjectRelease, error)
}

func NewMockProjectReleaseService() *MockProjectReleaseService {
	return &MockProjectReleaseService{}
}

func (m *MockProjectReleaseService) GetReleaseByIdWithDependencies(ctx context.Context, projectIdentifier string, releaseId string, userId string) (*models.ProjectRelease, error) {
	if m.GetReleaseByIdWithDependenciesFunc != nil {
		return m.GetReleaseByIdWithDependenciesFunc(ctx, projectIdentifier, releaseId, userId)
	}
	return nil, nil
}

func (m *MockProjectReleaseService) CreateRelease(ctx context.Context, projectIdentifier string, userId string, params CreateReleaseParams) (*models.ProjectRelease, error) {
	if m.CreateReleaseFunc != nil {
		return m.CreateReleaseFunc(ctx, projectIdentifier, userId, params)
	}
	return nil, nil
}

func (m *MockProjectReleaseService) GenerateProjectReleasePresignedPutUrl(ctx context.Context, projectIdentifier string, userId string, fileSize string) (string, error) {
	if m.GeneratePresignedPutUrlFunc != nil {
		return m.GeneratePresignedPutUrlFunc(ctx, projectIdentifier, userId, fileSize)
	}
	return "", nil
}

func (m *MockProjectReleaseService) GetReleasesByProjectId(ctx context.Context, id string, userId string) ([]models.ProjectRelease, error) {
	if m.GetReleasesByProjectIdFunc != nil {
		return m.GetReleasesByProjectIdFunc(ctx, id, userId)
	}
	return nil, nil
}

type MockLoaderVersionService struct {
	GetLoaderVersionByIdFunc          func(ctx context.Context, id string) (*models.LoaderVersion, error)
	GetLoaderVersionByGameVersionFunc func(ctx context.Context, gameVersion string) (*models.LoaderVersion, error)
	GetLoaderVersionByLabelFunc       func(ctx context.Context, label string) (*models.LoaderVersion, error)
	GetLoaderVersionsFunc             func(ctx context.Context) ([]models.LoaderVersion, error)
	CreateLoaderVersionFunc           func(ctx context.Context, params CreateLoaderVersionParams) error
}

func NewMockLoaderVersionService() *MockLoaderVersionService {
	return &MockLoaderVersionService{}
}

func (m *MockLoaderVersionService) GetLoaderVersionById(ctx context.Context, id string) (*models.LoaderVersion, error) {
	if m.GetLoaderVersionByIdFunc != nil {
		return m.GetLoaderVersionByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockLoaderVersionService) GetLoaderVersionByGameVersion(ctx context.Context, gameVersion string) (*models.LoaderVersion, error) {
	if m.GetLoaderVersionByGameVersionFunc != nil {
		return m.GetLoaderVersionByGameVersionFunc(ctx, gameVersion)
	}
	return nil, nil
}

func (m *MockLoaderVersionService) GetLoaderVersionByLabel(ctx context.Context, label string) (*models.LoaderVersion, error) {
	if m.GetLoaderVersionByLabelFunc != nil {
		return m.GetLoaderVersionByLabelFunc(ctx, label)
	}
	return nil, nil
}

func (m *MockLoaderVersionService) GetLoaderVersions(ctx context.Context) ([]models.LoaderVersion, error) {
	if m.GetLoaderVersionsFunc != nil {
		return m.GetLoaderVersionsFunc(ctx)
	}
	return nil, nil
}

func (m *MockLoaderVersionService) CreateLoaderVersion(ctx context.Context, params CreateLoaderVersionParams) error {
	if m.CreateLoaderVersionFunc != nil {
		return m.CreateLoaderVersionFunc(ctx, params)
	}
	return nil
}
