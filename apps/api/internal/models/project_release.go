package models

import "time"

type ProjectRelease struct {
	Id              string
	ProjectId       string
	Name            string
	Changelog       *string
	VersionNumber   string
	LoaderVersionId string
	LoaderVersion   LoaderVersion
	Downloads       int64
	FileUrl         string
	FileSize        int64
	FileHash        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	PublishedAt     *time.Time
	Dependencies    []ProjectReleaseDependency
}

type ProjectReleaseDependencyType string

const (
	ProjectReleaseDependencyTypeRequired ProjectReleaseDependencyType = "required"
	ProjectReleaseDependencyTypeOptional ProjectReleaseDependencyType = "optional"
)

type ProjectReleaseDependency struct {
	Id                  string
	ReleaseId           string
	DependencyProjectId string
	MinVersionNumber    *string
	Type                ProjectReleaseDependencyType
	CreatedAt           time.Time
}
