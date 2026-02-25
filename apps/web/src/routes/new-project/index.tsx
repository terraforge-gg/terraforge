import { createFileRoute, notFound } from "@tanstack/react-router";

export const Route = createFileRoute("/new-project/")({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (!context.session || !context.user) {
      throw notFound();
    }
  },
});

function RouteComponent() {
  return <div>Hello "/new-project/"!</div>;
}
