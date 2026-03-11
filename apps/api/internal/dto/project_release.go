package dto

import (
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type CreateProjectReleaseRequestDependency struct {
	ProjectId        string  `json:"projectId"`
	MinVersionNumber *string `json:"minVersionNumber"`
	Type             string  `json:"type" validate:"project_version_dependency_type"`
}

type CreateProjectReleaseRequest struct {
	Name            string                                  `json:"name" validate:"required,min=3,max=100"`
	VersionNumber   string                                  `json:"versionNumber" validate:"required"`
	Changelog       *string                                 `json:"changelog"`
	LoaderVersionId string                                  `json:"loaderVersionId" validate:"required"`
	FileUrl         string                                  `json:"fileUrl" validate:"required,file_url"`
	Dependencies    []CreateProjectReleaseRequestDependency `json:"dependencies"`
}

type ProjectReleaseDependencyResponse struct {
	Id               string    `json:"id"`
	ReleaseId        string    `json:"-"`
	ProjectId        string    `json:"projectId"`
	MinVersionNumber *string   `json:"minVersionNumber"`
	Type             string    `json:"type"`
	CreatedAt        time.Time `json:"createdAt"`
}

type ProjectReleaseResponse struct {
	Id            string                             `json:"id"`
	ProjectId     string                             `json:"projectId"`
	Name          string                             `json:"name"`
	Changelog     *string                            `json:"changelog,omitempty"`
	VersionNumber string                             `json:"versionNumber"`
	LoaderVersion LoaderVersionResponse              `json:"loaderVersion"`
	Downloads     int64                              `json:"downloads"`
	FileUrl       string                             `json:"fileUrl"`
	FileSize      int64                              `json:"fileSize"`
	FileHash      string                             `json:"fileHash"`
	CreatedAt     time.Time                          `json:"createdAt"`
	UpdatedAt     time.Time                          `json:"updatedAt"`
	PublishedAt   *time.Time                         `json:"publishedAt,omitempty"`
	Dependencies  []ProjectReleaseDependencyResponse `json:"dependencies"`
}

func MapToProjectReleaseDependencyResponse(d models.ProjectReleaseDependency) ProjectReleaseDependencyResponse {
	return ProjectReleaseDependencyResponse{
		Id:               d.Id,
		ReleaseId:        d.ReleaseId,
		ProjectId:        d.DependencyProjectId,
		MinVersionNumber: d.MinVersionNumber,
		Type:             string(d.Type),
		CreatedAt:        d.CreatedAt,
	}
}

func MapToProjectReleaseResponse(v models.ProjectRelease) ProjectReleaseResponse {
	deps := make([]ProjectReleaseDependencyResponse, len(v.Dependencies))

	for i, d := range v.Dependencies {
		deps[i] = MapToProjectReleaseDependencyResponse(d)
	}

	return ProjectReleaseResponse{
		Id:            v.Id,
		ProjectId:     v.ProjectId,
		Name:          v.Name,
		Changelog:     v.Changelog,
		VersionNumber: v.VersionNumber,
		LoaderVersion: MapToLoaderVersionResponse(v.LoaderVersion),
		Downloads:     v.Downloads,
		FileUrl:       v.FileUrl,
		FileSize:      v.FileSize,
		FileHash:      v.FileHash,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		PublishedAt:   v.PublishedAt,
		Dependencies:  deps,
	}
}
