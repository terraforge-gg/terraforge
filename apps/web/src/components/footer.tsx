import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/button";

const Footer = () => {
  return (
    <footer className="bg-muted h-41 flex items-center justify-between">
      <div className="text-muted-foreground w-full px-1 text-center text-xs leading-loose sm:text-sm">
        <a
          href="https://github.com/terraforge-gg/terraforge"
          className={cn(buttonVariants({ variant: "ghost" }))}
        >
          GitHub
        </a>
        <a
          href="https://status.terraforge.gg"
          className={cn(buttonVariants({ variant: "ghost" }))}
        >
          Status
        </a>
      </div>
    </footer>
  );
};

export default Footer;
