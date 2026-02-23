import { serve } from "@hono/node-server"
import { Hono } from "hono"
import { env } from "./env.js"
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { auth } from "./auth.js";

const app = new Hono()

app.use(logger())
app.use(
	"/api/auth/*",
	cors({
		origin: [env.FRONTEND_URL],
		allowHeaders: ["Content-Type", "Authorization"],
		allowMethods: ["POST", "GET", "OPTIONS"],
		exposeHeaders: ["Content-Length"],
		maxAge: 600,
		credentials: true,
	}),
);

app.get("/", (c) => {
  return c.json({
    env: env.APP_ENV,
    timestamp: new Date().toISOString(),
  });
});

app.get("/health", (c) => {
  return c.json({
    status: "ok",
    timestamp: new Date().toISOString(),
  });
});

app.on(["POST", "GET"], "/api/auth/**", (c) => auth.handler(c.req.raw));

serve({
	fetch: app.fetch,
	port: env.HOST_PORT,
});
