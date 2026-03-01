import type { components } from "./schema";

export type ProjectIdentifier = components["parameters"]["ProjectIdentifier"];
export type Project = components["schemas"]["Project"];

export type ProjectType = components["schemas"]["ProjectType"];
export const PROJECT_TYPES = ["mod"] as const satisfies ProjectType[];
export const PROJECT_TYPE = {
  MOD: "mod",
} as const satisfies Record<string, ProjectType>;

export type ProjectStatus = components["schemas"]["ProjectStatus"];
export const PROJECT_STATUSES = [
  "draft",
  "rejected",
  "approved",
  "banned",
] as const satisfies ProjectStatus[];
export const PROJECT_STATUS = {
  DRAFT: "draft",
  REJECTED: "rejected",
  APPROVED: "approved",
  BANNED: "banned",
} as const satisfies Record<string, ProjectStatus>;

export type ProjectMemberRole = components["schemas"]["ProjectMemberRole"];
export const PROJECT_MEMBER_ROLES = [
  "owner",
  "admin",
  "developer",
  "maintainer",
  "member",
] as const satisfies ProjectMemberRole[];
export const PROJECT_MEMBER_ROLE = {
  OWNER: "owner",
  ADMIN: "admin",
  DEVELOPER: "developer",
  MAINTAINER: "maintainer",
  MEMBER: "member",
} as const satisfies Record<string, ProjectMemberRole>;

export type ProjectMember = components["schemas"]["ProjectMember"];
export type ProjectSearch = components["schemas"]["ProjectSearch"];
