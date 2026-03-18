"use client";
import { ProjectMember } from "@/lib/api/types";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";

type ProjectMembersProps = {
  members: ProjectMember[];
};

const ProjectMembers = ({ members }: ProjectMembersProps) => {
  return (
    <Accordion type="single" collapsible defaultValue="members">
      <AccordionItem value="members">
        <AccordionTrigger className="tracking-widest">
          MAINTAINERS
        </AccordionTrigger>
        <AccordionContent>
          {members.map((x) => (
            <div key={x.userId} className="flex items-center gap-4">
              <Link href={`/user/${x.username}`}>{x.username}</Link>
              <Badge
                variant="outline"
                className="font-mono tracking-widest uppercase"
              >
                {x.role}
              </Badge>
            </div>
          ))}
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
};

export default ProjectMembers;
