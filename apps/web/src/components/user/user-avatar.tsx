import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

interface Props {
  avatar: string | undefined | null;
  fallback: string;
  className?: string;
}

const UserAvatar = ({ avatar, fallback, className }: Props) => {
  return (
    <Avatar className={cn("h-9 w-9", className)}>
      <AvatarImage src={avatar as string | undefined} />
      <AvatarFallback>{fallback}</AvatarFallback>
    </Avatar>
  );
};

export const UserAvatarSkeleton = () => {
  return <Skeleton className="h-9 w-9 rounded-full" />;
};

export default UserAvatar;
