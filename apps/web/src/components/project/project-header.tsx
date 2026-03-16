"use client";
import Image from "next/image";
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
import Link from "../link";

type ProjecInfo = Pick<
  Project,
  "name" | "slug" | "summary" | "iconUrl" | "downloads"
>;
type ProjectReleaseInfo = Pick<ProjectRelease, "versionNumber" | "fileUrl">;

type ProjectHeaderProps = {
  latestVersionDownloadLink?: string;
  project: ProjecInfo;
  latestRelease?: ProjectReleaseInfo;
  showSettings?: boolean;
};

const ProjectHeader = ({
  project,
  latestRelease,
  showSettings,
}: ProjectHeaderProps) => {
  const { name, summary, iconUrl, downloads } = project;

  const iconSize = "w-24 h-24";

  return (
    <header className="sticky top-12 flex w-full flex-col gap-6 bg-background pt-4 pb-10">
      <Item className="px-0">
        <>
          <ItemMedia variant="image" className={iconSize}>
            {iconUrl ? (
              <Image src={iconUrl} alt={name} fill />
            ) : (
              <BoxIcon className={iconSize} />
            )}
          </ItemMedia>
          <ItemContent>
            <ItemTitle>
              <h1 className="font-mono text-3xl">{name}</h1>
            </ItemTitle>
            <ItemDescription className="text-sm text-muted-foreground">
              <div className="flex flex-col">
                <span>{summary}</span>
                <div className="flex items-center gap-2">
                  <div className="flex items-center gap-2">
                    <DownloadIcon size={16} />
                    {downloads}
                  </div>
                  {/* <Separator orientation="vertical" />
                  <Badge variant="outline">DLC</Badge> */}
                </div>
              </div>
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
      <div>
        <div className="flex items-center gap-4">
          <Link
            href={`/mod/${project.slug}`}
            className="border-b border-transparent font-mono transition ease-in-out hover:border-chart-2 hover:text-chart-2"
            exact
            activeProps={{
              className: "text-chart-1 border-chart-1",
            }}
          >
            description
          </Link>
          {showSettings && (
            <Link
              href={`/mod/${project.slug}/settings`}
              className="border-b border-transparent font-mono transition ease-in-out hover:border-chart-2 hover:text-chart-2"
              exact
              activeProps={{
                className: "text-chart-1 border-chart-1",
              }}
            >
              settings
            </Link>
          )}
        </div>
        <Separator />
      </div>
    </header>
  );
};

export default ProjectHeader;
