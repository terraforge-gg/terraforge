-- +goose Up
-- +goose StatementBegin
CREATE TYPE "project_type" AS ENUM ('mod');

CREATE TYPE "project_status" AS ENUM ('draft', 'rejected', 'approved', 'banned');

CREATE TABLE "project" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "slug" TEXT UNIQUE NOT NULL,
    "summary" TEXT,
    "description" TEXT,
    "iconUrl" TEXT,
    "downloads" INTEGER DEFAULT 0 NOT NULL,
    "type" project_type DEFAULT 'mod' NOT NULL,
    "status" project_status DEFAULT 'draft' NOT NULL,
    "createdAt" TIMESTAMP DEFAULT now() NOT NULL,
    "updatedAt" TIMESTAMP DEFAULT now() NOT NULL,
    "deletedAt" TIMESTAMP DEFAULT NULL,
    "userId" TEXT REFERENCES "user" ("id") ON DELETE CASCADE
);

CREATE INDEX "project_userId_idx" ON "project"("userId");
CREATE INDEX "project_type_idx" ON "project"("type");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "project_type_idx";

DROP INDEX "project_userId_idx";

DROP TABLE "project";

DROP TYPE "project_status";

DROP TYPE "project_type";
-- +goose StatementEnd