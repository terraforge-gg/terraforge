import { LogOutIcon, PlusIcon, User2Icon } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { Button, buttonVariants } from "@/components/ui/button";
import { signOut, useSession } from "@/lib/auth-client";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import UserAvatar, { UserAvatarSkeleton } from "@/components/user/user-avatar";
import { useState } from "react";

const Navbar = () => {
  const { data: session, isPending } = useSession();
  const [open, setOpen] = useState(false);

  return (
    <header className="light:bg-gray-100/60 inset-x-0 top-0 z-10 py-4 hidden sm:block">
      <div className="flex h-full items-center justify-between gap-2">
        <div className="flex items-center gap-8 text-xl">
          <Link to="/">
            <span className="text-foreground font-extrabold">terraforge</span>
          </Link>
          <Link
            to="/"
            search={{ query: undefined, page: undefined, perPage: undefined }}
            className="text-foreground hover:decoration-primary text-lg font-semibold decoration-2 underline-offset-4 hover:underline"
            activeProps={{ className: "decoration-primary underline" }}
          >
            Home
          </Link>
        </div>
        <div className="flex items-center justify-end space-x-8">
          {!isPending && session && (
            <DropdownMenu open={open} onOpenChange={setOpen}>
              <DropdownMenuTrigger asChild>
                <Button>
                  <PlusIcon />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-40" align="start">
                <DropdownMenuItem asChild>
                  <Link
                    className="hover:cursor-pointer"
                    to="/new-project"
                    onClick={() => setOpen(false)}
                  >
                    Create Project
                  </Link>
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          )}
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
                  <Link to="/">
                    <DropdownMenuItem className="hover:bg-accent hover:cursor-pointer">
                      <User2Icon className="mr-2 h-4 w-4" />
                      <span>Profile</span>
                    </DropdownMenuItem>
                  </Link>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="hover:bg-accent hover:cursor-pointer"
                    onClick={() => signOut()}
                  >
                    <LogOutIcon className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ) : (
            <div className="flex gap-4">
              <Link
                to="/sign-in"
                className={buttonVariants({ variant: "outline" })}
              >
                Sign In
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Navbar;
