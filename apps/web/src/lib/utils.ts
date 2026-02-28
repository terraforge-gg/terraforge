import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getChangedFields<T extends object>(
  current: T,
  defaults: T,
): Partial<T> {
  return Object.fromEntries(
    Object.entries(current).filter(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      ([key, value]) => value !== (defaults as any)[key],
    ),
  ) as Partial<T>;
}
