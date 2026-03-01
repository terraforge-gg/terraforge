import z from "zod";
import { PROJECT_TYPE, PROJECT_TYPES } from "../../types";

export const searchProjectSchema = z
  .object({
    query: z
      .string()
      .transform((val) => (val === "" ? undefined : val))
      .optional(),
    type: z.enum(PROJECT_TYPES).default(PROJECT_TYPE.MOD),
    page: z.coerce.number().optional(),
    perPage: z.coerce.number().optional(),
  })
  .optional();

export type SearchProjectSchema = z.infer<typeof searchProjectSchema>;

export type SearchProjectParams =
  | {
      query?: string;
      page?: number;
      perPage?: number;
    }
  | undefined;
