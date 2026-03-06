package dto

import (
	"time"

	"github.com/terraforge-gg/terraforge/internal/models"
)

type LoaderVersionResponse struct {
	Id              string    `json:"id"`
	GameVersion     string    `json:"gameVersion"`
	InternalVersion string    `json:"internalVersion"`
	Status          string    `json:"status"`
	IsLegacy        bool      `json:"isLegacy"`
	ReleasedAt      time.Time `json:"releasedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func MapToLoaderVersionResponse(lv models.LoaderVersion) LoaderVersionResponse {
	return LoaderVersionResponse{
		Id:              lv.Id,
		GameVersion:     lv.GameVersion,
		InternalVersion: lv.InternalVersion,
		Status:          string(lv.Status),
		IsLegacy:        lv.IsLegacy,
		ReleasedAt:      lv.ReleasedAt,
		CreatedAt:       lv.CreatedAt,
		UpdatedAt:       lv.UpdatedAt,
	}
}
