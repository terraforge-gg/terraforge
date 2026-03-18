"use client";
import { ProjectRelease } from "@/lib/api/types";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { DownloadIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import CreateProjectReleaseDialog from "./create-release-dialog";
import { useQuery } from "@tanstack/react-query";
import { getProjectReleasesQueryOptions } from "@/lib/api/query-options/project-release";
import { Spinner } from "../ui/spinner";

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
  const { data: releases, isPending } = useQuery(
    getProjectReleasesQueryOptions(
      { projectSlug: projectSlug },
      {
        placeholderData: (previousData) => previousData,
        initialData: initialReleases,
      },
    ),
  );

  return (
    <Accordion type="single" collapsible defaultValue="releases">
      <AccordionItem value="releases">
        <AccordionTrigger className="tracking-widest">
          RELEASES
        </AccordionTrigger>
        <AccordionContent className="flex h-auto flex-col gap-2">
          {showCreateRelease && (
            <CreateProjectReleaseDialog
              projectId={projectId}
              projectSlug={projectSlug}
            />
          )}
          {isPending && (
            <div className="flex justify-center">
              <Spinner />
            </div>
          )}
          {releases && !isPending && releases.length > 0 ? (
            releases.map((x, i) => (
              <div
                key={x.id}
                className={cn(
                  "flex justify-between rounded-md p-2 hover:bg-accent",
                  i === 0 && "rounded-md bg-accent p-2",
                )}
              >
                <div className="flex flex-col items-center">
                  <span className="font-mono text-lg">{x.versionNumber}</span>
                  {i === 0 && (
                    <span className="font-mono text-xs text-chart-1">
                      LATEST
                    </span>
                  )}
                </div>
                <div className="flex items-center justify-center gap-2">
                  <DownloadIcon size={16} />
                  {x.downloads}
                </div>
              </div>
            ))
          ) : (
            <div className="text-muted-foreground">NO RELEASES FOUND</div>
          )}
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

export default ProjectReleases;
