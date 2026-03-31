"use client";
import Image from "next/image";
import { BoxIcon } from "lucide-react";
import { Item, ItemContent, ItemMedia, ItemTitle } from "@/components/ui/item";

type UserHeaderProps = {
  username?: string;
  image?: string;
};

const iconSize = "w-24 h-24";

const UserHeader = ({ username, image }: UserHeaderProps) => {
  return (
    <header className="flex w-full flex-col gap-6 pt-4 pb-6">
      <Item className="px-0">
        <>
          <ItemMedia variant="image" className={iconSize}>
            {image ? (
              <Image src={image} alt={username ?? "user image"} fill />
            ) : (
              <BoxIcon className={iconSize} />
            )}
          </ItemMedia>
          <ItemContent>
            <ItemTitle>
              <h1 className="font-mono text-3xl">{username}</h1>
            </ItemTitle>
          </ItemContent>
        </>
      </Item>
    </header>
  );
};

export default UserHeader;
