"use client";
import { BoxIcon, DownloadIcon } from "lucide-react";
import type { Project, ProjectRelease } from "@/lib/api/types";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import { Separator } from "@/components/ui/separator";
import { Button } from "../ui/button";
import { toast } from "sonner";
import { Badge } from "../ui/badge";

type ProjecInfo = Pick<Project, "name" | "summary" | "iconUrl" | "downloads">;
type ProjectReleaseInfo = Pick<ProjectRelease, "versionNumber" | "fileUrl">;

type ProjectHeaderProps = {
  latestVersionDownloadLink?: string;
  project: ProjecInfo;
  latestRelease?: ProjectReleaseInfo;
};

const ProjectHeader = ({ project, latestRelease }: ProjectHeaderProps) => {
  const { name, summary, iconUrl, downloads } = project;

  return (
    <header className="sticky top-12 flex w-full flex-col gap-6 bg-background pt-10 pb-10">
      <Item>
        <>
          <ItemMedia variant="image" className="h-12 w-12">
            <BoxIcon className="h-12 w-12" />
          </ItemMedia>
          <ItemContent>
            <ItemTitle>
              <h1 className="font-mono text-3xl">{name}</h1>
              <span className="font-mono text-lg text-muted-foreground">
                {latestRelease && `v${latestRelease.versionNumber}`}
              </span>
            </ItemTitle>
            <ItemDescription className="flex gap-2 font-mono">
              <span className="flex gap-2">
                <DownloadIcon size={20} /> {}
                {project.downloads}
              </span>
              <Separator orientation="vertical" />
              <Badge>Content</Badge>
              <Badge>QoL</Badge>
              <Badge>World Gen</Badge>
            </ItemDescription>
          </ItemContent>
          <ItemContent className="flex-none text-center">
            <ItemDescription>
              <Button
                size="lg"
                onClick={() => toast.info("not yet implemented")}
              >
                <DownloadIcon />
                download
              </Button>
            </ItemDescription>
          </ItemContent>
        </>
      </Item>
    </header>
  );
};

export default ProjectHeader;
