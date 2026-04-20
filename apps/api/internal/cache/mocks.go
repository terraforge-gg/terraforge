package cache

import (
	"context"
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type MockProjectCache struct {
	GetProjectFunc        func(ctx context.Context, identifier string) (*models.Project, error)
	SetProjectFunc        func(ctx context.Context, project *models.Project, ttl time.Duration) error
	DeleteProjectFunc     func(ctx context.Context, id string) error
	GetProjectMembersFunc func(ctx context.Context, identifier string) ([]models.ProjectMember, error)
	SetProjectMembersFunc func(ctx context.Context, project *models.Project, members []models.ProjectMember, ttl time.Duration) error
}

func NewMockProjectCache() *MockProjectCache {
	return &MockProjectCache{}
}

func (m *MockProjectCache) GetProject(ctx context.Context, identifier string) (*models.Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(ctx, identifier)
	}
	return nil, ErrCacheMiss
}

func (m *MockProjectCache) SetProject(ctx context.Context, project *models.Project, ttl time.Duration) error {
	if m.SetProjectFunc != nil {
		return m.SetProjectFunc(ctx, project, ttl)
	}
	return nil
}

func (m *MockProjectCache) DeleteProject(ctx context.Context, id string) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectCache) GetProjectMembers(ctx context.Context, identifier string) ([]models.ProjectMember, error) {
	if m.GetProjectMembersFunc != nil {
		return m.GetProjectMembersFunc(ctx, identifier)
	}
	return nil, ErrCacheMiss
}

func (m *MockProjectCache) SetProjectMembers(ctx context.Context, project *models.Project, members []models.ProjectMember, ttl time.Duration) error {
	if m.SetProjectMembersFunc != nil {
		return m.SetProjectMembersFunc(ctx, project, members, ttl)
	}
	return nil
}
