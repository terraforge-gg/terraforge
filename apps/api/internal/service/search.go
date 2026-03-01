package service

import (
	"context"
	"log/slog"

	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/repository"
)

type SearchService interface {
	SearchProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) ([]models.Project, int64)
	Health(ctx context.Context) error
}

type searchService struct {
	logger     *slog.Logger
	searchRepo repository.SearchRepository
}

func NewSearchService(logger *slog.Logger, searchRepo repository.SearchRepository) SearchService {
	return &searchService{logger: logger, searchRepo: searchRepo}
}

func (s *searchService) SearchProjects(ctx context.Context, query string, projectType string, limit int64, offset int64) ([]models.Project, int64) {
	result, err := s.searchRepo.FindProjects(ctx, query, projectType, limit, offset)

	if err != nil {
		panic(err)
	}

	projects := make([]models.Project, len(result.Projects))

	for i, p := range result.Projects {
		projects[i] = models.Project{
			Id:          p.Id,
			Name:        p.Name,
			Slug:        p.Slug,
			Summary:     p.Summary,
			IconUrl:     p.IconUrl,
			Description: p.Description,
			Downloads:   p.Downloads,
			Type:        models.ProjectType(p.Type),
			Status:      models.ProjectStatus(p.Status),
			UpdatedAt:   p.UpdatedAt,
			CreatedAt:   p.CreatedAt,
			UserId:      p.UserId,
		}
	}

	return projects, result.TotalHits
}

func (s *searchService) Health(ctx context.Context) error {
	err := s.searchRepo.Health(ctx)

	if err != nil {
		return err
	}

	return nil
}
