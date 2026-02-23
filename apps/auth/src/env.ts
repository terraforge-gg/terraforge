import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";
import dotenv from "dotenv";

dotenv.config();

export const env = createEnv({
  server: {
    NODE_ENV: z
      .enum(["development", "production"])
      .default("development"),
    APP_ENV: z.string(),
    HOST_PORT: z.coerce.number(),
    DATABASE_URL: z.url(),
    BETTER_AUTH_SECRET: z.string().nonempty(),
    BETTER_AUTH_URL: z.url(),
    FRONTEND_URL: z.url(),
    DISCORD_CLIENT_ID: z.string(),
    DISCORD_CLIENT_SECRET: z.string(),
  },
  runtimeEnv: process.env,
  emptyStringAsUndefined: true,
});