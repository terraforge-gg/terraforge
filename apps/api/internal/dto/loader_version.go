package dto

import (
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type LoaderVersionResponse struct {
	Id           string    `json:"id"`
	GameVersion  string    `json:"gameVersion"`
	VersionLabel string    `json:"versionLabel"`
	BuildType    string    `json:"buildType"`
	ReleasedAt   time.Time `json:"releasedAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func MapToLoaderVersionResponse(lv models.LoaderVersion) LoaderVersionResponse {
	return LoaderVersionResponse{
		Id:           lv.Id,
		GameVersion:  lv.GameVersion,
		VersionLabel: lv.VersionLabel,
		BuildType:    string(lv.BuildType),
		ReleasedAt:   lv.ReleasedAt,
		UpdatedAt:    lv.UpdatedAt,
	}
}
