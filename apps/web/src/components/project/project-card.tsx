import Image from "next/image";
import { BoxIcon, Download, RefreshCcw } from "lucide-react";
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
import Link from "next/link";

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
      <ItemMedia variant="image" className="w-24 h-24">
        {iconUrl ? (
          <Image src={iconUrl} alt={name} fill />
        ) : (
          <BoxIcon className="w-24 h-24" />
        )}
      </ItemMedia>
      <ItemContent>
        <ItemTitle className="font-bold text-lg md:text-2xl line-clamp-1">
          <Link href={`/mods/${slug}`} className="hover:underline">
            {name}
          </Link>
          {ownerName && " - "}
          {ownerName && (
            <Link
              href="/"
              className="font-bold text-base text-primary hover:underline"
            >
              {ownerName}
            </Link>
          )}
        </ItemTitle>
        <ItemDescription className="text-ellipsis break-all">
          {summary}
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex flex-row sm:flex-col items-start">
        <div className="flex items-center gap-2">
          <Download size={16} />
          <span className="font-bold">
            {downloads.toLocaleString("en-US")}
          </span>{" "}
          downloads
        </div>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger className="hover:cursor-default">
              <div className="flex items-center gap-2">
                <RefreshCcw size={16} />
                {formatDistance(updatedAt, new Date(), { addSuffix: true })}
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <div>{format(updatedAt, "MMMM d, yyyy 'at' h:m a")}</div>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </ItemActions>
    </Item>
  );
};

export default ProjectCard;
