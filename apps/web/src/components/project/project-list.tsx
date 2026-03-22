"use client";
import ProjectCard from "./project-card";
import type { Project } from "@/lib/api/types";
import { ItemGroup } from "@/components/ui/item";
import { Empty, EmptyHeader, EmptyTitle } from "@/components/ui/empty";
import { Skeleton } from "@/components/ui/skeleton";

interface ProjectListProps {
  projects?: Project[];
  loading?: boolean;
}

const ProjectList = ({ projects, loading }: ProjectListProps) => {
  if (loading) {
    return <Skeleton className="h-96 w-full" />;
  }

  return (
    <div className="w-full py-4">
      {projects && projects.length ? (
        <ItemGroup className="gap-4">
          {projects.map((x) => (
            <ProjectCard
              key={x.id}
              name={x.name}
              slug={x.slug}
              summary={x.summary}
              iconUrl={x.iconUrl}
              downloads={x.downloads}
              type={x.type}
              updatedAt={x.updatedAt}
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
