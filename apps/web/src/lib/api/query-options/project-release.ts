import { queryOptions } from "@tanstack/react-query";
import type { UseQueryOptions } from "@tanstack/react-query";
import type { ProjectRelease } from "@/lib/api/types";
import api from "@/lib/api/api";
import {
  GetProjectReleasePresignedPutUrlParams,
  GetProjectReleasesParams,
} from "../models/project/create-release";

export const projectReleasePresignedPutUrlQueryOptions = (
  params: GetProjectReleasePresignedPutUrlParams,
  options?: Omit<
    UseQueryOptions<string, Error, string>,
    "queryKey" | "queryFn"
  >,
) => {
  return queryOptions({
    ...options,
    queryKey: ["project-release", "presigned-put", { ...params }],
    queryFn: () => api.project.getProjectReleasePresignedPutUrl(params),
    retry: false,
    refetchOnWindowFocus: false,
  });
};

export const getProjectReleasesQueryOptions = (
  params: GetProjectReleasesParams,
  options?: Omit<
    UseQueryOptions<Array<ProjectRelease>, Error, Array<ProjectRelease>>,
    "queryKey" | "queryFn"
  >,
) => {
  return queryOptions({
    ...options,
    queryKey: ["project-releases", { ...params }],
    queryFn: () => api.project.releases(params),
    refetchOnWindowFocus: false,
  });
};
