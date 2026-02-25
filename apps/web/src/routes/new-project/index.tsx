import { createFileRoute, notFound } from "@tanstack/react-router";
import CreateProjectForm from "@/components/project/create-project-form";

export const Route = createFileRoute("/new-project/")({
  component: RouteComponent,
  beforeLoad: ({ context }) => {
    if (!context.session || !context.user) {
      throw notFound();
    }
  },
});

function RouteComponent() {
  return (
    <div className="flex justify-center">
      <CreateProjectForm />
    </div>
  );
}
