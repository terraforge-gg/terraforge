-- +goose Up
-- +goose StatementBegin
CREATE TYPE "loader_version_build_type" AS ENUM ('stable', 'preview', 'legacy');

CREATE TABLE "loader_version" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "gameVersion" TEXT NOT NULL,
    "versionLabel" TEXT NOT NULL,
    "buildType" loader_version_build_type NOT NULL,
    "releasedAt" TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMP DEFAULT now() NOT NULL,
    UNIQUE("gameVersion", "versionLabel", "buildType")
);

CREATE INDEX "loader_version_buildType_idx" ON "loader_version"("buildType");

CREATE INDEX "loader_version_versionLabel_idx" ON "loader_version"("versionLabel");

CREATE INDEX "loader_version_gameVersion_idx" ON "loader_version"("gameVersion");

CREATE INDEX "loader_version_releasedAt_idx" ON "loader_version"("releasedAt" DESC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX "loader_version_releasedAt_idx";

DROP INDEX "loader_version_gameVersion_idx";

DROP INDEX "loader_version_versionLabel_idx";

DROP INDEX "loader_version_buildType_idx";

DROP TABLE "loader_version";

DROP TYPE loader_version_build_type;
-- +goose StatementEnd