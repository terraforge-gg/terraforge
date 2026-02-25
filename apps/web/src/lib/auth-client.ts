import { createAuthClient } from "better-auth/react";
import { usernameClient } from "better-auth/client/plugins";
import { env } from "@/env/client";

export const { signIn, signUp, signOut, useSession, getSession } =
  createAuthClient({
    baseURL: env.VITE_AUTH_URL,
    plugins: [usernameClient()],
  });
