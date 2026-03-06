package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

var loaderVersions = []CreateLoaderVersionParams{
	{Id: utils.GenerateUUID(), GameVersion: "1.3.5", InternalVersion: "0.11.8.10", Status: string(models.LoaderVersionStatusStable), IsLegacy: true, ReleasedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{Id: utils.GenerateUUID(), GameVersion: "1.4.3", InternalVersion: "2022.09.47.88", Status: string(models.LoaderVersionStatusStable), IsLegacy: true, ReleasedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{Id: utils.GenerateUUID(), GameVersion: "1.4.4", InternalVersion: "2026.02.2.2", Status: string(models.LoaderVersionStatusStable), IsLegacy: false, ReleasedAt: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func SeedLoaderVersions(logger *slog.Logger, loaderVersionService LoaderVersionService) error {
	ctx := context.Background()
	for _, lv := range loaderVersions {
		_, err := loaderVersionService.GetLoaderVersionByGameVersion(ctx, lv.GameVersion)
		if err == nil {
			logger.Info("loader version already exists, skipping", "version", lv.GameVersion)
			continue
		}

		if err != errors.ErrLoaderVersionNotFound {
			return err
		}

		if err := loaderVersionService.CreateLoaderVersion(ctx, lv); err != nil {
			return err
		}

		logger.Info("seeded loader version", "version", lv.GameVersion)
	}

	return nil
}
