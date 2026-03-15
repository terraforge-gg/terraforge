package models

import "time"

type LoaderVersionBuildType string

const (
	LoaderVersionStatusStable  LoaderVersionBuildType = "stable"
	LoaderVersionStatusPreview LoaderVersionBuildType = "preview"
	LoaderVersionStatusLegacy  LoaderVersionBuildType = "legacy"
)

type LoaderVersion struct {
	Id           string
	GameVersion  string
	VersionLabel string
	BuildType    LoaderVersionBuildType
	ReleasedAt   time.Time
	UpdatedAt    time.Time
}
