import { queryOptions, UseQueryOptions } from "@tanstack/react-query";
import { ProjectSearch } from "../types";
import { SearchProjectParams } from "../models/project/search";
import api from "../api";

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
    queryFn: () => api.project.search(params),
  });
};
