import Image from "next/image";
import { BoxIcon, DownloadIcon } from "lucide-react";
import type { Project } from "@/lib/api/types";
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import { Separator } from "@/components/ui/separator";

type ProjecInfo = Pick<Project, "name" | "summary" | "iconUrl" | "downloads">;

type ProjectHeaderProps = {
  latestVersionDownloadLink?: string;
} & ProjecInfo;

const ProjectHeader = ({
  name,
  summary,
  iconUrl,
  downloads,
}: ProjectHeaderProps) => {
  return (
    <>
      <div className="flex w-full flex-col gap-6">
        <Item>
          <ItemMedia variant="image" className="w-24 h-24">
            {iconUrl ? (
              <Image src={iconUrl} alt={name} fill />
            ) : (
              <BoxIcon className="w-24 h-24" />
            )}
          </ItemMedia>
          <ItemContent>
            <ItemTitle className="font-bold text-lg md:text-4xl">
              {name}
            </ItemTitle>
            <ItemDescription>
              <span className="flex gap-2 items-center">
                <DownloadIcon size={16} />
                <span className="text-end text-sm">{downloads}</span>
              </span>
            </ItemDescription>
            <ItemDescription>{summary}</ItemDescription>
          </ItemContent>
          <ItemActions className="items-center hidden sm:flex"></ItemActions>
        </Item>
      </div>
      <Separator />
    </>
  );
};

export default ProjectHeader;
