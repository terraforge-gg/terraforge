import { betterAuth } from "better-auth";
import { v7 as uuidv7 } from 'uuid';
import { jwt, openAPI, username } from "better-auth/plugins"
import { env } from "./env.js";
import { Pool } from "pg";

export const auth = betterAuth({
  baseURL: env.BETTER_AUTH_URL,
  basePath: "/api/auth",
  trustedOrigins: [env.FRONTEND_URL],
  emailAndPassword: {
    enabled: true
  },
  socialProviders: {
    discord: {
      clientId: env.DISCORD_CLIENT_ID,
      clientSecret: env.DISCORD_CLIENT_SECRET,
      mapProfileToUser: (profile) => {
        return {
          username: profile.username,
        };
      },
      redirectURI: env.FRONTEND_URL
    },
  },
  session: {
    cookieCache: {
      enabled: false,
    }
  },
  advanced: {
    cookies: {
      session_token: {
        name: "auth_token",
        attributes: {
          httpOnly: true,
          secure: true,
          domain: env.NODE_ENV === "development" ? "localhost" : ".terraforge.gg",
          sameSite: "None",
        }
      }
    },
    database: {
      generateId: () => uuidv7(),

    }
  },
  plugins: [
    username(),
    jwt(),
    openAPI(),
  ],
  database: new Pool({
    connectionString: env.DATABASE_URL,
  }),
});