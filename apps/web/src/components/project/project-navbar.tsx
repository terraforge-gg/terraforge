import { cn } from "@/lib/utils";

import { buttonVariants } from "@/components/ui/button";
import Link from "../link";

interface ProjectNavbarProps {
  slug: string;
  className?: string;
  showSettings?: boolean;
}

const ProjectNavbar = ({
  slug,
  className,
  showSettings,
}: ProjectNavbarProps) => {
  return (
    <div className={cn("flex gap-2 py-4 flex-wrap items-center", className)}>
      <Link
        href={`/mods/${slug}`}
        className={cn(buttonVariants({ variant: "link" }), "text-foreground")}
        activeProps={{ className: "underline" }}
      >
        Description
      </Link>
      <Link
        href={`/mods/${slug}/images`}
        className={cn(buttonVariants({ variant: "link" }), "text-foreground")}
        activeProps={{ className: "underline" }}
      >
        Images
      </Link>
      <Link
        href={`/mods/${slug}/versions`}
        className={cn(buttonVariants({ variant: "link" }), "text-foreground")}
        activeProps={{ className: "underline" }}
      >
        Versions
      </Link>
      {showSettings && (
        <Link
          href={`/mods/${slug}/settings`}
          className={cn(buttonVariants({ variant: "link" }), "text-foreground")}
          activeProps={{ className: "underline" }}
        >
          Settings
        </Link>
      )}
    </div>
  );
};

export default ProjectNavbar;
