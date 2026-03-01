package meilisearch

import (
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type ProjectDocument struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Summary     *string   `json:"summary"`
	Description *string   `json:"description"`
	IconUrl     *string   `json:"iconUrl"`
	Downloads   int64     `json:"downloads"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UserId      string    `json:"userId"`
}

func ProjectToDocument(project *models.Project) *ProjectDocument {
	return &ProjectDocument{
		Id:          project.Id,
		Name:        project.Name,
		Slug:        project.Slug,
		Summary:     project.Summary,
		Description: project.Description,
		IconUrl:     project.IconUrl,
		Downloads:   project.Downloads,
		Type:        string(project.Type),
		Status:      string(project.Status),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		UserId:      project.UserId,
	}
}

type ProjectSearchResult struct {
	Projects  []ProjectDocument
	TotalHits int64
}
