-- +goose Up
-- +goose StatementBegin
CREATE TYPE "project_type" AS ENUM ('mod');

CREATE TYPE "project_status" AS ENUM ('draft', 'rejected', 'approved', 'banned');

CREATE TABLE "project" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
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

CREATE UNIQUE INDEX "project_slug_unique_idx" ON "project"("slug") WHERE "deletedAt" IS NULL;

CREATE INDEX "project_userId_idx" ON "project"("userId");
CREATE INDEX "project_type_idx" ON "project"("type");

CREATE VIEW "active_project" AS SELECT * FROM "project" WHERE "deletedAt" IS NULL;

CREATE TYPE "project_member_role" AS ENUM ('owner', 'admin', 'developer', 'maintainer', 'member');

CREATE TABLE "project_member" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "projectId" TEXT NOT NULL REFERENCES "project" ("id") ON DELETE CASCADE,
    "userId" TEXT REFERENCES "user" ("id") ON DELETE CASCADE,
    "role" project_member_role DEFAULT 'member' NOT NULL,
    "createdAt" TIMESTAMP DEFAULT now() NOT NULL,
    UNIQUE("projectId", "userId")
);

CREATE INDEX "project_member_projectId_userId_idx" ON "project_member"("projectId", "userId");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX "project_member_projectId_userId_idx";

DROP TABLE "project_member";

DROP TYPE "project_member_role";

DROP INDEX "project_type_idx";

DROP INDEX "project_userId_idx";

DROP TABLE "project";

DROP TYPE "project_status";

DROP TYPE "project_type";
-- +goose StatementEnd