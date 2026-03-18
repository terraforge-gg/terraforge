import { queryOptions, UseQueryOptions } from "@tanstack/react-query";
import api from "../api";
import { LoaderVersion } from "../types";

export const loaderVersionsQueryOptions = (
  options?: Omit<
    UseQueryOptions<LoaderVersion[], Error, LoaderVersion[]>,
    "queryKey" | "queryFn"
  >,
) => {
  return queryOptions({
    ...options,
    queryKey: ["loader-versions"],
    queryFn: () => api.loaderVersion.list(),
    refetchOnWindowFocus: false,
    staleTime: 1000 * 60 * 60 * 24,
  });
};
