import { cn } from "@/lib/utils";
import { buttonVariants } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

const Footer = () => {
  return (
    <>
      <Separator />
      <footer className="flex h-24 items-center">
        <div className="flex w-full justify-center gap-4 px-1 text-center text-xs leading-loose text-muted-foreground sm:text-sm">
          <a
            href="https://github.com/terraforge-gg/terraforge"
            className={cn(buttonVariants({ variant: "ghost" }))}
          >
            github
          </a>
          <a
            href="https://status.terraforge.gg"
            className={cn(buttonVariants({ variant: "ghost" }))}
          >
            status
          </a>
        </div>
      </footer>
    </>
  );
};

export default Footer;
