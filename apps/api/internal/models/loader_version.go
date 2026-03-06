package models

import "time"

type LoaderVersionStatus string

const (
	LoaderVersionStatusPreview LoaderVersionStatus = "preview"
	LoaderVersionStatusStable  LoaderVersionStatus = "stable"
)

type LoaderVersion struct {
	Id              string
	GameVersion     string
	InternalVersion string
	Status          LoaderVersionStatus
	IsLegacy        bool
	ReleasedAt      time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
