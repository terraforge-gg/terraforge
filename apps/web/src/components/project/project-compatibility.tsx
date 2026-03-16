"use client";
import { ProjectRelease } from "@/lib/api/types";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

type ProjectCompatibilityProps = {
  release: ProjectRelease;
};

const ProjectCompatibility = ({ release }: ProjectCompatibilityProps) => {
  return (
    <Accordion type="single" collapsible defaultValue="members">
      <AccordionItem value="members">
        <AccordionTrigger className="tracking-widest">
          COMPATIBILITY
        </AccordionTrigger>
        <AccordionContent>
          <span>{release.loaderVersion.gameVersion}</span>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

export default ProjectCompatibility;
