"use client";

import { cn } from "@/lib/utils";
import NextLink, { LinkProps as NextLinkProps } from "next/link";
import { usePathname } from "next/navigation";
import { AnchorHTMLAttributes, ReactNode } from "react";

type ActiveProps = Omit<AnchorHTMLAttributes<HTMLAnchorElement>, "href">;

interface LinkProps
  extends NextLinkProps, Omit<AnchorHTMLAttributes<HTMLAnchorElement>, "href"> {
  children?: ReactNode;
  /** Props applied when the link href matches the current pathname */
  activeProps?: ActiveProps;
  /** Match the pathname exactly (default: true) */
  exact?: boolean;
}

export function Link({
  href,
  children,
  activeProps,
  exact = true,
  className,
  ...props
}: LinkProps) {
  const pathname = usePathname();
  const hrefString = typeof href === "string" ? href : (href.pathname ?? "");

  const isActive = exact
    ? pathname === hrefString
    : pathname.startsWith(hrefString);

  return (
    <NextLink
      href={href}
      className={cn(className, isActive && activeProps?.className)}
      {...props}
    >
      {children}
    </NextLink>
  );
}

export default Link;
