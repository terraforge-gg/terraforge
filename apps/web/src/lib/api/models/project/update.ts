import z from "zod";
import {
  PROJECT_NAME_MAX_LENGTH,
  PROJECT_NAME_MIN_LENGTH,
  slugSchema,
} from "./create";
import type { ProjectIdentifier } from "../../types";

export const updateProjectSchema = z.object({
  name: z
    .string()
    .min(PROJECT_NAME_MIN_LENGTH, {
      error: `Name must be longer than ${PROJECT_NAME_MIN_LENGTH} characters`,
    })
    .max(PROJECT_NAME_MAX_LENGTH, {
      error: `Name must be less than ${PROJECT_NAME_MAX_LENGTH} characters`,
    })
    .optional(),
  summary: z.string().optional(),
  slug: slugSchema.optional(),
  description: z.string().optional(),
  iconUrl: z.url().optional(),
});

export type UpdateProjectSchema = z.infer<typeof updateProjectSchema>;

export type UpdateProjectParams = {
  projectIdentifier: ProjectIdentifier;
  values: UpdateProjectSchema;
};
