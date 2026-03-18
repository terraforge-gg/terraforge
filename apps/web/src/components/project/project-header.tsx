"use client";
import Image from "next/image";
import { BoxIcon, DownloadIcon } from "lucide-react";
import type { Project } from "@/lib/api/types";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import { Separator } from "@/components/ui/separator";
import Link from "../link";
import DownloadReleaseButton from "./download-release-button";
import { getProjectReleasesQueryOptions } from "@/lib/api/query-options/project-release";
import { useQuery } from "@tanstack/react-query";
import { Spinner } from "@/components/ui/spinner";

type ProjecInfo = Pick<
  Project,
  "name" | "slug" | "summary" | "iconUrl" | "downloads"
>;

type ProjectHeaderProps = {
  project: ProjecInfo;
  showSettings?: boolean;
};

const ProjectHeader = ({ project, showSettings }: ProjectHeaderProps) => {
  const { name, summary, iconUrl, downloads } = project;

  const { data, isPending } = useQuery(
    getProjectReleasesQueryOptions({
      projectSlug: project.slug,
    }),
  );

  const latestRelease = data?.[0];

  const iconSize = "w-24 h-24";

  return (
    <header className="flex w-full flex-col gap-6 pt-4 pb-6">
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
                    <>
                      {isPending ? (
                        <Spinner />
                      ) : latestRelease ? (
                        <>
                          <DownloadIcon size={16} />
                          {downloads}{" "}
                        </>
                      ) : null}
                    </>
                  </div>
                  {/* <Separator orientation="vertical" />
                  <Badge variant="outline">DLC</Badge> */}
                </div>
              </div>
            </ItemDescription>
          </ItemContent>
          {latestRelease && (
            <ItemContent className="flex-none text-center">
              <ItemDescription>
                <DownloadReleaseButton
                  fileUrl={latestRelease.fileUrl}
                  size="icon-lg"
                />
              </ItemDescription>
            </ItemContent>
          )}
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
