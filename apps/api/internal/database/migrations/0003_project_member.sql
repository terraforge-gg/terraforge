-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd