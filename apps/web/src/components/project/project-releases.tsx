"use client";
import { ProjectRelease } from "@/lib/api/types";
import CreateProjectReleaseDialog from "./create-release-dialog";
import { useQuery } from "@tanstack/react-query";
import { getProjectReleasesQueryOptions } from "@/lib/api/query-options/project-release";

type ProjectReleasesProps = {
  projectId: string;
  projectSlug: string;
  initialReleases: ProjectRelease[];
  showCreateRelease?: boolean;
};

const ProjectReleases = ({
  projectId,
  projectSlug,
  initialReleases,
  showCreateRelease,
}: ProjectReleasesProps) => {
  const { data: releases } = useQuery(
    getProjectReleasesQueryOptions(
      { projectSlug: projectSlug },
      {
        placeholderData: (previousData) => previousData,
        initialData: initialReleases,
      },
    ),
  );

  return (
    <div className="flex flex-col gap-4">
      <div className="font-mono text-sm">LATEST RELEASES</div>
      <div className="flex flex-col gap-2">
        {showCreateRelease && (
          <CreateProjectReleaseDialog
            projectId={projectId}
            projectSlug={projectSlug}
          />
        )}
        {releases && releases.length > 0 ? (
          <div className="divide-y divide-border">
            {releases?.map((x) => (
              <div
                key={x.id}
                className="flex items-center justify-between py-1"
              >
                <div className="font-mono text-sm">{x.versionNumber}</div>
                <div className="font-mono text-sm">{x.publishedAt}</div>
              </div>
            ))}
          </div>
        ) : (
          <div>NO RELEASES</div>
        )}
      </div>
    </div>
  );
};

export default ProjectReleases;
