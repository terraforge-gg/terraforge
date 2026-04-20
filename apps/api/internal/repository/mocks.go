package repository

import (
	"context"

	"github.com/terraforge-gg/terraforge/internal/lib/meilisearch"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type MockSearchRepository struct {
	Projects  []meilisearch.ProjectDocument
	TotalHits int64
	Err       error
}

func NewMockSearchRepository() *MockSearchRepository {
	return &MockSearchRepository{}
}

func (m *MockSearchRepository) IndexProject(ctx context.Context, project *models.Project) error {
	return nil
}

func (m *MockSearchRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	return nil
}

func (m *MockSearchRepository) DeleteProject(ctx context.Context, projectId string) error {
	return nil
}

func (m *MockSearchRepository) FindProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) (*meilisearch.ProjectSearchResult, error) {
	return &meilisearch.ProjectSearchResult{
		Projects:  m.Projects,
		TotalHits: m.TotalHits,
	}, m.Err
}

func (m *MockSearchRepository) Health(ctx context.Context) error {
	return m.Err
}

func (m *MockSearchRepository) EnsureProjectIndexExists() {}
