import createClient from "openapi-fetch";
import type {Middleware} from "openapi-fetch";
import type { paths } from "./schema";
import { env } from "@/env/client";
import { token } from "../auth-client";

export const client = createClient<paths>({
  baseUrl: env.VITE_API_URL + "/" + env.VITE_API_VERSION,
});

const PROTECTED_ROUTES = ["/projects"];
const PROTECTED_METHODS = ["POST", "PATCH", "PUT", "DELETE"];

const auth: Middleware = {
  async onRequest({ request, schemaPath }) {
    if (
      PROTECTED_METHODS.includes(request.method) &&
      PROTECTED_ROUTES.some((pathname) => schemaPath.startsWith(pathname))
    ) {
      const { data, error } = await token();

      if (error) {
        throw new Error(error.message);
      }

      request.headers.set("Authorization", `Bearer ${data.token}`);
    }

    return request;
  },
};

client.use(auth);
