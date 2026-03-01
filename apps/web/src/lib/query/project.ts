import { queryOptions } from "@tanstack/react-query";
import type { UseQueryOptions } from "@tanstack/react-query";
import apiService from "@/lib/api/service";
import type { Project, ProjectMember, ProjectSearch } from "../api/types";
import { SearchProjectParams } from "../api/models/project/search";

export const projectQueryOptions = (
  identifier: string,
  options?: Omit<
    UseQueryOptions<Project | null, Error, Project | null>,
    "queryKey" | "queryFn"
  >,
) =>
  queryOptions({
    ...options,
    queryKey: ["project", { identifier: identifier }],
    queryFn: () => apiService.project.identifier(identifier),
    refetchOnWindowFocus: false,
  });

export const projectMembersQueryOptions = (
  identifier: string,
  options?: Omit<
    UseQueryOptions<
      Array<ProjectMember> | null,
      Error,
      Array<ProjectMember> | null
    >,
    "queryKey" | "queryFn"
  >,
) =>
  queryOptions({
    ...options,
    queryKey: ["project-members", { identifier: identifier }],
    queryFn: () => apiService.project.members(identifier),
    refetchOnWindowFocus: false,
  });

export const projectSearchQueryOptions = (
  params?: SearchProjectParams,
  options?: Omit<
    UseQueryOptions<ProjectSearch, Error, ProjectSearch>,
    "queryKey" | "queryFn"
  >,
) => {
  return queryOptions({
    ...options,
    queryKey: ["project-search", params],
    queryFn: () => apiService.project.search(params),
  });
};
