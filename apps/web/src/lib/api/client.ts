import createClient from "openapi-fetch";
import type { Middleware } from "openapi-fetch";
import type { paths } from "./schema";
import { env } from "@/env";
import { token } from "../auth-client";

export const client = createClient<paths>({
  baseUrl: env.NEXT_PUBLIC_API_URL + "/" + env.NEXT_PUBLIC_API_VERSION,
  credentials: "include",
});

const auth: Middleware = {
  async onRequest({ request }) {
    let cookie: string | undefined;

    if (typeof window === "undefined") {
      try {
        const { cookies } = await import("next/headers");
        const cookieStore = await cookies();
        cookie = cookieStore.toString();
      } catch {
        cookie = undefined;
      }
    }

    const options = cookie
      ? {
          fetchOptions: {
            headers: {
              cookie: cookie,
            },
          },
        }
      : undefined;

    const { data, error } = await token(options);

    if (data && !error) {
      request.headers.set("Authorization", `Bearer ${data.token}`);
    }

    return request;
  },
};

client.use(auth);
