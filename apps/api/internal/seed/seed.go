package seed

import (
	"context"
	"log/slog"
	"time"

	"github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/utils"
)

var loaderVersions = []service.CreateLoaderVersionParams{
	{Id: utils.NewUUID(), GameVersion: "1.3.5", VersionLabel: "0.11.8.10", BuildType: string(models.LoaderVersionStatusLegacy), ReleasedAt: time.Now(), UpdatedAt: time.Now()},
	{Id: utils.NewUUID(), GameVersion: "1.4.3", VersionLabel: "2022.09.47.88", BuildType: string(models.LoaderVersionStatusLegacy), ReleasedAt: time.Now(), UpdatedAt: time.Now()},
	{Id: utils.NewUUID(), GameVersion: "1.4.4", VersionLabel: "2026.02.2.3", BuildType: string(models.LoaderVersionStatusStable), ReleasedAt: time.Now(), UpdatedAt: time.Now()},
	{Id: utils.NewUUID(), GameVersion: "1.4.4", VersionLabel: "2026.01.3.2", BuildType: string(models.LoaderVersionStatusStable), ReleasedAt: time.Now(), UpdatedAt: time.Now()},
}

func SeedLoaderVersions(logger *slog.Logger, loaderVersionService service.LoaderVersionService) error {
	ctx := context.Background()
	for _, lv := range loaderVersions {
		_, err := loaderVersionService.GetLoaderVersionByLabel(ctx, lv.GameVersion)
		if err == nil {
			logger.Info("loader version already exists, skipping", "game version", lv.GameVersion, "version label", lv.VersionLabel)
			continue
		}

		if err != errors.ErrLoaderVersionNotFound {
			return err
		}

		if err := loaderVersionService.CreateLoaderVersion(ctx, lv); err != nil {
			return err
		}

		logger.Info("seeded loader version", "game version", lv.GameVersion, "version label", lv.VersionLabel)
	}

	return nil
}
