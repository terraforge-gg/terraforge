"use client";
import Image from "next/image";
import { BoxIcon, Download, RefreshCcw } from "lucide-react";
import millify from "millify";
import { format, formatDistance } from "date-fns";
import type { Project } from "@/lib/api/types";
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import Link from "../link";
type MinimalProject = Omit<
  Project,
  "id" | "description" | "status" | "createdAt" | "organisationId"
>;

type ProjectCardProps = {
  username?: string;
} & MinimalProject;

const iconSize = "w-24 h-24";

const ProjectCard = ({
  slug,
  name,
  iconUrl,
  summary,
  downloads,
  updatedAt,
  username,
}: ProjectCardProps) => {
  return (
    <Item variant="muted" className="items-start p-4">
      <ItemMedia variant="image" className={iconSize}>
        {iconUrl ? (
          <Image src={iconUrl} alt={name} fill />
        ) : (
          <BoxIcon className={iconSize} />
        )}
      </ItemMedia>
      <ItemContent>
        <ItemTitle className="line-clamp-1 text-lg font-bold break-normal text-ellipsis md:text-2xl">
          <Link href={`/mod/${slug}`} className="hover:underline">
            {name}
          </Link>
          {username && " - "}
          {username && (
            <Link
              href={`/user/${username}`}
              className="text-base font-bold text-primary hover:underline"
            >
              {username}
            </Link>
          )}
        </ItemTitle>
        <ItemDescription className="break-all text-ellipsis">
          {summary}
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex flex-col items-start">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger className="hover:cursor-default">
              <div className="flex items-center gap-2 text-xs md:text-lg">
                <Download size={20} />
                <div className="flex gap-1">
                  <span className="font-bold">{millify(downloads)}</span>{" "}
                  downloads
                </div>
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <span>{`${downloads.toLocaleString("en-US")} downloads`}</span>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>

        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger className="hover:cursor-default">
              <div className="flex items-center gap-2 text-xs md:text-lg">
                <RefreshCcw size={20} />
                {formatDistance(updatedAt, new Date(), { addSuffix: true })}
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <span>{format(updatedAt, "MMMM d, yyyy 'at' h:m a")}</span>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </ItemActions>
    </Item>
  );
};

export default ProjectCard;
