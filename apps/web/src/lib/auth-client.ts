import { createAuthClient } from "better-auth/react";
import { jwtClient, usernameClient } from "better-auth/client/plugins";
import { env } from "@/env/client";

export const { signIn, signUp, signOut, useSession, getSession, token } =
  createAuthClient({
    baseURL: env.VITE_AUTH_URL,
    plugins: [usernameClient(), jwtClient()],
  });
