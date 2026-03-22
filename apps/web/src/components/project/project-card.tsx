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

interface ProjectCardProps extends MinimalProject {
  ownerName?: string;
}

const ProjectCard = ({
  slug,
  name,
  iconUrl,
  summary,
  downloads,
  updatedAt,
  ownerName,
}: ProjectCardProps) => {
  return (
    <Item variant="muted" className="items-start p-4">
      <ItemMedia variant="image" className="h-24 w-24">
        {iconUrl ? (
          <img src={iconUrl} alt={name} />
        ) : (
          <BoxIcon className="h-24 w-24" />
        )}
      </ItemMedia>
      <ItemContent>
        <ItemTitle className="line-clamp-1 text-lg font-bold md:text-2xl">
          <Link href={`/mod/${slug}`} className="hover:underline">
            {name}
          </Link>
          {ownerName && " - "}
          {ownerName && (
            <Link
              href={`/user/${ownerName}`}
              className="text-base font-bold text-primary hover:underline"
            >
              {ownerName}
            </Link>
          )}
        </ItemTitle>
        <ItemDescription className="break-normal text-ellipsis">
          {summary}
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex flex-row items-start sm:flex-col">
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger className="hover:cursor-default">
              <div className="flex items-center gap-2 text-lg">
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
              <div className="flex items-center gap-2 text-lg">
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
