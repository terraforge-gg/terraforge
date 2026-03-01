import ProjectCard from "./project-card";
import type { Project } from "@/lib/api/types";
import { ItemGroup } from "@/components/ui/item";
import { Empty, EmptyHeader, EmptyTitle } from "@/components/ui/empty";
import { Skeleton } from "@/components/ui/skeleton";

interface ProjectListProps {
  projects: Array<Project>;
  loading?: boolean;
}

const ProjectList = ({ projects, loading }: ProjectListProps) => {
  if (loading) {
    return <Skeleton className="w-full h-96" />;
  }

  return (
    <div className="py-4 w-full">
      {projects.length ? (
        <ItemGroup className="gap-4">
          {projects.map((x) => (
            <ProjectCard
              key={x.id}
              name={x.name}
              slug={x.slug}
              summary={x.summary}
              iconUrl={x.iconUrl}
              downloads={x.downloads}
              updatedAt={x.updatedAt}
              type={x.type}
            />
          ))}
        </ItemGroup>
      ) : (
        <Empty>
          <EmptyHeader>
            <EmptyTitle>No projects found.</EmptyTitle>
          </EmptyHeader>
        </Empty>
      )}
    </div>
  );
};

export default ProjectList;
