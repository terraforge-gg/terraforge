"use client";

import { signOut, useSession } from "@/lib/auth/client";
import { useRouter } from "next/navigation";
import Link from "./link";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOutIcon, User2Icon } from "lucide-react";
import UserAvatar, { UserAvatarSkeleton } from "@/components/user/user-avatar";
import SignInDialog from "./user/sign-in-dialog";
import CreateProjectDialog from "@/components/project/create-project-dialog";
import { Badge } from "./ui/badge";

const Navbar = () => {
  const { data: session, isPending } = useSession();
  const router = useRouter();

  return (
    <header className="light:bg-gray-100/60 sticky inset-x-0 top-0 z-10 border-b py-2 backdrop-blur-sm sm:block">
      <div className="container mx-auto flex h-full max-w-6xl items-center justify-between gap-2 px-4">
        <div className="flex items-center gap-2 text-lg">
          <Link href="/">
            <span className="text-foreground">terraforge</span>
          </Link>
          <Badge className="font-mono" variant="outline">
            alpha
          </Badge>
        </div>
        <div className="flex items-center justify-end space-x-8">
          {!isPending && session && <CreateProjectDialog />}
          {isPending ? (
            <UserAvatarSkeleton />
          ) : session?.user ? (
            <div className="flex gap-10">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <div className="relative rounded-full hover:cursor-pointer">
                    <UserAvatar
                      avatar={session.user.image}
                      fallback={session.user.username?.[0] || ""}
                    />
                  </div>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuLabel>My Account</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <Link href={`/user/${session.user.username}`}>
                    <DropdownMenuItem className="hover:cursor-pointer hover:bg-accent">
                      <User2Icon className="mr-2 h-4 w-4" />
                      <span>Profile</span>
                    </DropdownMenuItem>
                  </Link>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="hover:cursor-pointer hover:bg-accent"
                    onClick={() => {
                      signOut();
                      router.push("/");
                    }}
                  >
                    <LogOutIcon className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ) : (
            <div className="flex gap-4">
              <SignInDialog />
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Navbar;
