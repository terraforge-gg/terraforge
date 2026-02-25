import z from "zod";

export const PROJECT_NAME_MIN_LENGTH = 3;
export const PROJECT_NAME_MAX_LENGTH = 100;

export const slugSchema = z
  .string()
  .min(1, "Slug cannot be empty")
  .max(100, "Slug must be less than 100 characters")
  .regex(
    /^[a-zA-Z0-9_-]+$/,
    "Slug must contain only alphanumeric characters, hyphens, and underscores",
  )
  .regex(/^[a-zA-Z0-9]/, "Slug must start with an alphanumeric character")
  .regex(/[a-zA-Z0-9]$/, "Slug must end with an alphanumeric character");

export const createProjectSchema = z.object({
  type: z.enum(["mod"]),
  name: z
    .string()
    .min(PROJECT_NAME_MIN_LENGTH, {
      error: `Name must be longer than ${PROJECT_NAME_MIN_LENGTH} characters`,
    })
    .max(PROJECT_NAME_MAX_LENGTH, {
      error: `Name must be less than ${PROJECT_NAME_MAX_LENGTH} characters`,
    }),
  slug: slugSchema,
  summary: z.string().optional(),
  organisationId: z.uuidv7().optional(),
});

export type CreateProjectSchema = z.infer<typeof createProjectSchema>;
