-- +goose Up
-- +goose StatementBegin
CREATE TABLE "project_release" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "projectId" TEXT NOT NULL REFERENCES "project" ("id") ON DELETE CASCADE,
    "name" TEXT NOT NULL,
    "changelog" TEXT,
    "versionNumber" TEXT NOT NULL,
    "loaderVersionId" TEXT NOT NULL REFERENCES "loader_version" ("id") ON DELETE RESTRICT,
    "downloads" INTEGER DEFAULT 0 NOT NULL,
    "fileUrl" TEXT NOT NULL,
    "fileSize" INTEGER NOT NULL,
    "fileHash" TEXT NOT NULL,
    "createdAt" TIMESTAMP DEFAULT now() NOT NULL,
    "updatedAt" TIMESTAMP DEFAULT now() NOT NULL,
    "publishedAt" TIMESTAMP,
    UNIQUE("projectId", "versionNumber"),
    CHECK ("versionNumber" ~ '^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-[a-zA-Z0-9.]+)?$')
);

CREATE INDEX "project_release_projectId_idx" ON "project_release"("projectId");

CREATE TYPE "project_release_dependency_type" AS ENUM ('required', 'optional');

CREATE TABLE "project_release_dependency" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "releaseId" TEXT NOT NULL REFERENCES "project_release" ("id") ON DELETE CASCADE,
    "dependencyProjectId" TEXT NOT NULL REFERENCES "project" ("id") ON DELETE CASCADE,
    "minVersionNumber" TEXT,
    "type" project_release_dependency_type DEFAULT 'required' NOT NULL,
    "createdAt" TIMESTAMP DEFAULT now() NOT NULL,
    UNIQUE("releaseId", "dependencyProjectId")
);

CREATE INDEX "project_release_dependency_releaseId_idx" ON "project_release_dependency"("releaseId");

CREATE INDEX "project_release_dependency_dependencyProjectId_idx" ON "project_release_dependency"("dependencyProjectId");
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE "project_release_dependency";

DROP INDEX "project_release_projectId_idx";

DROP TABLE "project_release";

DROP TYPE "project_release_dependency_type";
-- +goose StatementEnd