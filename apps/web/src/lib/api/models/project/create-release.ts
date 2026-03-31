import z from "zod";
import { env } from "@/env";
import { ProjectIdentifier } from "../../types";

export type GetProjectReleasesParams = {
  projectSlug: ProjectIdentifier;
};

export type GetProjectReleasePresignedPutUrlParams = {
  projectId: ProjectIdentifier;
  fileSize: number;
};

export const PROJECT_RELEASE_NAME_MIN_LENGTH = 3;
export const PROJECT_RELEASE_NAME_MAX_LENGTH = 100;
export const PROJECT_RELEASE_MAX_DEPENDENCIES = 16;

const projectReleaseDependencySchema = z.object({
  type: z.enum(["required", "optional"]),
  projectId: z.uuidv7(),
  minVersionNumber: z.string().optional(),
});

const semverSchema = z
  .string()
  .regex(
    /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-[a-zA-Z0-9.]+)?$/,
    "Must be valid semver (e.g., 1.2.3)",
  );

export const createProjectReleaseSchema = z.object({
  name: z
    .string()
    .min(PROJECT_RELEASE_NAME_MIN_LENGTH, {
      error: `Name must be longer than ${PROJECT_RELEASE_NAME_MIN_LENGTH} characters`,
    })
    .max(PROJECT_RELEASE_NAME_MAX_LENGTH, {
      error: `Name must be less than ${PROJECT_RELEASE_NAME_MAX_LENGTH} characters`,
    }),
  versionNumber: semverSchema,
  changelog: z.string().optional(),
  loaderVersionId: z.uuidv7({ error: "Invalid loader version" }),
  fileUrl: z.url().refine((url) => {
    try {
      return new URL(url).origin === env.NEXT_PUBLIC_CDN_URL;
    } catch {
      return false;
    }
  }),
  dependencies: z
    .array(projectReleaseDependencySchema)
    .max(PROJECT_RELEASE_MAX_DEPENDENCIES)
    .optional(),
});

export type CreateProjectReleaseSchema = z.infer<
  typeof createProjectReleaseSchema
>;
export type CreateProjectReleaseParams = {
  projectId: string;
  values: CreateProjectReleaseSchema;
};
