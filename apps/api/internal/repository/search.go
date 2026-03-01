package repository

import (
	"context"
	"encoding/json"
	"log/slog"

	msearch "github.com/meilisearch/meilisearch-go"
	"github.com/terraforge-gg/terraforge/internal/lib/meilisearch"
	"github.com/terraforge-gg/terraforge/internal/models"
)

type SearchRepository interface {
	IndexProject(ctx context.Context, project *models.Project) error
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, projectId string) error
	FindProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) (*meilisearch.ProjectSearchResult, error)
	Health(ctx context.Context) error
	EnsureProjectIndexExists()
}

type meiliSearchRepository struct {
	logger      *slog.Logger
	meiliSearch *meilisearch.MeiliSearchClient
}

func NewMeiliSearchRepository(logger *slog.Logger, meiliSearchClient *meilisearch.MeiliSearchClient) SearchRepository {
	repo := &meiliSearchRepository{
		logger:      logger,
		meiliSearch: meiliSearchClient,
	}

	repo.EnsureProjectIndexExists()

	return repo
}

const PROJECTS_INDEX = "projects"

func (s *meiliSearchRepository) IndexProject(ctx context.Context, project *models.Project) error {
	doc := meilisearch.ProjectToDocument(project)
	index := s.meiliSearch.Client.Index(PROJECTS_INDEX)
	_, err := index.AddDocuments([]meilisearch.ProjectDocument{*doc}, &msearch.DocumentOptions{PrimaryKey: msearch.StringPtr("id")})
	return err
}

func (s *meiliSearchRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	doc := meilisearch.ProjectToDocument(project)
	index := s.meiliSearch.Client.Index(PROJECTS_INDEX)
	_, err := index.UpdateDocuments([]meilisearch.ProjectDocument{*doc}, &msearch.DocumentOptions{PrimaryKey: msearch.StringPtr("id")})
	return err
}

func (s *meiliSearchRepository) DeleteProject(ctx context.Context, projectId string) error {
	index := s.meiliSearch.Client.Index(PROJECTS_INDEX)
	_, err := index.DeleteDocument(projectId, nil)
	return err
}

func (s *meiliSearchRepository) FindProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) (*meilisearch.ProjectSearchResult, error) {
	index := s.meiliSearch.Client.Index(PROJECTS_INDEX)
	request := &msearch.SearchRequest{
		Limit:  limit,
		Offset: offset,
		Filter: "type = '" + projectType + "'",
	}

	res, err := index.Search(query, request)

	if err != nil {
		return nil, err
	}

	var projects []meilisearch.ProjectDocument
	for _, hit := range res.Hits {
		var project meilisearch.ProjectDocument

		hitJSON, err := json.Marshal(hit)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(hitJSON, &project); err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	result := &meilisearch.ProjectSearchResult{
		Projects:  projects,
		TotalHits: res.EstimatedTotalHits,
	}

	return result, nil
}

func (s *meiliSearchRepository) Health(ctx context.Context) error {
	_, err := s.meiliSearch.Client.Health()

	if err != nil {
		return err
	}

	return nil
}

func (s *meiliSearchRepository) EnsureProjectIndexExists() {
	_, err := s.meiliSearch.Client.GetIndex(PROJECTS_INDEX)
	if err != nil {

		_, err = s.meiliSearch.Client.CreateIndex(&msearch.IndexConfig{
			Uid:        PROJECTS_INDEX,
			PrimaryKey: "id",
		})
		if err != nil {
			s.logger.Error("Failed to create '"+PROJECTS_INDEX+"' index", "Error:", err)
		}

		s.logger.Info(PROJECTS_INDEX + " index created")
	}

	index := s.meiliSearch.Client.Index(PROJECTS_INDEX)
	_, err = index.UpdateSettings(&msearch.Settings{
		SearchableAttributes: []string{"name", "slug", "summary", "description"},
		FilterableAttributes: []string{"type", "downloads", "updatedAt"},
	})

	if err != nil {
		s.logger.Error("Failed to update projects index", "Error:", err)
	}

	s.logger.Info(PROJECTS_INDEX + " index setting updated")
}
