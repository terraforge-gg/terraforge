-- +goose Up
-- +goose StatementBegin
CREATE TYPE "loader_version_status" AS ENUM ('preview', 'stable');

CREATE TABLE "loader_version" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "gameVersion" TEXT NOT NULL,
    "internalVersion" TEXT NOT NULL,
    "status" loader_version_status NOT NULL,
    "isLegacy" BOOLEAN DEFAULT FALSE NOT NULL,
    "releasedAt" TIMESTAMP NOT NULL,
    "createdAt" TIMESTAMP NOT NULL,
    "updatedAt" TIMESTAMP NOT NULL
);

CREATE INDEX "loader_version_status_idx" ON "loader_version"("status");

CREATE INDEX "loader_version_gameVersion_idx" ON "loader_version"("gameVersion");

CREATE INDEX "loader_version_releasedAt_idx" ON "loader_version"("releasedAt" DESC);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX "loader_version_releasedAt_idx";

DROP INDEX "loader_version_status_idx";

DROP TABLE "loader_version";

DROP TYPE loader_version_status;
-- +goose StatementEnd