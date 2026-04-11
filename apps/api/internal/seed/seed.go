package seed

import (
	"context"
	"log/slog"
	"strings"

	"github.com/terraforge-gg/terraforge/internal/errors"
	"github.com/terraforge-gg/terraforge/internal/github"
	"github.com/terraforge-gg/terraforge/internal/models"
	"github.com/terraforge-gg/terraforge/internal/service"
	"github.com/terraforge-gg/terraforge/internal/utils"
	"github.com/terraforge-gg/terraforge/internal/validation"
)

func SeedLoaderVersions(logger *slog.Logger, loaderVersionService service.LoaderVersionService) error {
	ctx := context.Background()

	logger.Info("fetching tModLoader releases from GitHub")
	releases, err := github.FetchTModLoaderReleases(ctx)
	if err != nil {
		logger.Error("failed to fetch tModLoader releases from GitHub", "error", err)
		return err
	}

	logger.Info("fetched tModLoader releases from GitHub", "count", len(releases))

	for _, release := range releases {
		gameVersion := extractGameVersion(release.Name)
		if gameVersion == "" {
			logger.Warn("could not extract game version from release", "tag", release.TagName, "name", release.Name)
			continue
		}

		versionLabel := strings.TrimPrefix(release.TagName, "v")

		buildType := determineBuildType(release.Prerelease, release.Name)

		params := service.CreateLoaderVersionParams{
			Id:           utils.NewUUID(),
			GameVersion:  gameVersion,
			VersionLabel: versionLabel,
			BuildType:    string(buildType),
			ReleasedAt:   release.PublishedAt,
			UpdatedAt:    release.PublishedAt,
		}

		_, err := loaderVersionService.GetLoaderVersionByLabel(ctx, versionLabel)
		if err == nil {
			logger.Info("loader version already exists, skipping", "game version", gameVersion, "version label", versionLabel)
			continue
		}

		if err != errors.ErrLoaderVersionNotFound {
			return err
		}

		if err := loaderVersionService.CreateLoaderVersion(ctx, params); err != nil {
			return err
		}

		logger.Info("seeded loader version", "game version", gameVersion, "version label", versionLabel, "build type", buildType)
	}

	return nil
}

func extractGameVersion(name string) string {
	matches := validation.SemVerValidator.FindStringSubmatch(name)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func determineBuildType(prerelease bool, name string) models.LoaderVersionBuildType {
	if strings.Contains(name, "refs/heads/stable") {
		return models.LoaderVersionStatusStable
	}

	if strings.Contains(name, "refs/heads/preview") {
		return models.LoaderVersionStatusPreview
	}

	if prerelease {
		return models.LoaderVersionStatusPreview
	}

	return models.LoaderVersionStatusStable
}
