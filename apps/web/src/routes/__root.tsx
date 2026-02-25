import {
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from "@tanstack/react-router";
import { TanStackRouterDevtoolsPanel } from "@tanstack/react-router-devtools";
import { TanStackDevtools } from "@tanstack/react-devtools";

import TanStackQueryProvider from "@/integrations/tanstack-query/root-provider";

import TanStackQueryDevtools from "@/integrations/tanstack-query/devtools";

import appCss from "@/styles.css?url";

import type { QueryClient } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";
import Navbar from "@/components/navbar";
import Footer from "@/components/footer";
import { createServerFn } from "@tanstack/react-start";
import { getRequestHeaders } from "@tanstack/react-start/server";
import { getSession } from "@/lib/auth-client";
import requestLogger from "@/middleware/request-logger";

interface MyRouterContext {
  queryClient: QueryClient;
}

const getServerSession = createServerFn({ method: "GET" }).handler(async () => {
  const headers = getRequestHeaders();
  const response = await getSession({
    fetchOptions: {
      timeout: 10_000,
      headers: {
        cookie: headers.get("cookie") ?? "",
      },
    },
  });

  return response;
});

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: "utf-8",
      },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      {
        title: "terraforge",
      },
    ],
    links: [
      {
        rel: "stylesheet",
        href: appCss,
      },
    ],
  }),
  shellComponent: RootDocument,
  notFoundComponent: () => <div>not found</div>,
  errorComponent: () => (
    <div className="flex justify-center min-h-screen">
      <div className="mt-20">Something went wrong</div>
    </div>
  ),
  beforeLoad: async () => {
    try {
      const { data } = await getServerSession();

      return { ...data };
    } catch {
      return {};
    }
  },
  server: {
    middleware: [requestLogger],
  },
});

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="dark">
      <head>
        <HeadContent />
      </head>
      <body>
        <TanStackQueryProvider>
          <main className="container mx-auto max-w-7xl px-4 min-h-screen">
            <Navbar />
            <div className="py-12">{children}</div>
          </main>
          <Toaster />
          <Footer />
          <TanStackDevtools
            config={{
              position: "bottom-right",
            }}
            plugins={[
              {
                name: "Tanstack Router",
                render: <TanStackRouterDevtoolsPanel />,
              },
              TanStackQueryDevtools,
            ]}
          />
        </TanStackQueryProvider>
        <Scripts />
      </body>
    </html>
  );
}
