import { DownloadIcon, Loader2Icon } from "lucide-react";
import { useState } from "react";
import type { buttonVariants } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import type { VariantProps } from "class-variance-authority";

type DownloadReleaseButtonProps = {
  text?: string;
  fileUrl: string;
} & VariantProps<typeof buttonVariants>;

const DownloadReleaseButton = ({
  text,
  fileUrl,
  ...props
}: DownloadReleaseButtonProps) => {
  const [downloading, setDownloading] = useState(false);

  const onClick = () => {
    setDownloading(true);
    const link = document.createElement("a");
    link.href = fileUrl;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    setDownloading(false);
  };

  return (
    <Button {...props} onClick={onClick} disabled={downloading}>
      <div className="flex items-center justify-center gap-2">
        {downloading ? (
          <Loader2Icon className="animate-spin" />
        ) : (
          <DownloadIcon />
        )}
        {text}
      </div>
    </Button>
  );
};

export default DownloadReleaseButton;
