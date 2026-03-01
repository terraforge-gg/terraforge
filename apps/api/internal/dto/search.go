package dto

import "github.com/terraforge-gg/terraforge/internal/models"

type ProjectSearchResponse struct {
	Data      []ProjectResponse `json:"data"`
	TotalHits int64             `json:"totalHits"`
	Limit     int64             `json:"limit"`
	Offset    int64             `json:"offset"`
}

func ProjectToProjectSearchResponse(projects []models.Project, totalHits int64, limit int64, offset int64) ProjectSearchResponse {
	data := make([]ProjectResponse, len(projects))

	for i, p := range projects {
		data[i] = ProjectToProjectResponse(p)
	}

	return ProjectSearchResponse{
		Data:      data,
		TotalHits: totalHits,
		Limit:     limit,
		Offset:    offset,
	}
}
