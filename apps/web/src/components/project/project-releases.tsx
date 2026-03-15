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

type ProjectReleasesProps = {
  releases: ProjectRelease[];
};

const ProjectReleases = ({ releases }: ProjectReleasesProps) => {
  return (
    <Accordion type="single" collapsible>
      <AccordionItem value="releases">
        <AccordionTrigger>releases</AccordionTrigger>
        <AccordionContent>
          {releases.map((x, i) => (
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
                  <span className="font-mono text-xs text-chart-1">LATEST</span>
                )}
              </div>
              <div className="flex items-center justify-center gap-2">
                <DownloadIcon size={16} />
                {x.downloads}
              </div>
            </div>
          ))}
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

export default ProjectReleases;
