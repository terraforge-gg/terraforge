"use client";
import { ProjectMember } from "@/lib/api/types";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import UserAvatar from "../user/user-avatar";

type ProjectMembersProps = {
  members: ProjectMember[];
};

const ProjectMembers = ({ members }: ProjectMembersProps) => {
  return (
    <div className="flex flex-col gap-4">
      <div className="font-mono text-sm">TEAM MEMBERS</div>
      <div className="flex flex-col gap-2">
        {members.map((x) => (
          <div key={x.userId} className="flex items-center gap-2">
            <UserAvatar
              avatar={x.image}
              fallback={x.username[0]}
              className="h-6 w-6"
            />
            <Link
              href={`/user/${x.username}`}
              className="font-mono text-sm no-underline underline-offset-2 hover:underline"
            >
              ~{x.username}
            </Link>
            <Badge variant="outline" className="ml-auto font-mono">
              {x.role}
            </Badge>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ProjectMembers;
