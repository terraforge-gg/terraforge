"use client";
import ProjectList from "@/components/project/project-list";
import ProjectSearch from "@/components/project/project-search";
import { projectSearchQueryOptions } from "@/lib/api/query-options/project";
import { ProjectSearch as ProjectSearchResult } from "@/lib/api/types";
import { useQuery } from "@tanstack/react-query";
import { parseAsInteger, useQueryState } from "nuqs";
import { useDebounce } from "use-debounce";

type SearchInnerProps = {
  initialData?: ProjectSearchResult;
};

const Inner = ({ initialData }: SearchInnerProps) => {
  const [page, setPage] = useQueryState("page", parseAsInteger.withDefault(0));
  const [perPage, setPerPage] = useQueryState(
    "perPage",
    parseAsInteger.withDefault(10),
  );
  const [_query, setQuery] = useQueryState("query");
  const [query] = useDebounce(_query, 500);

  const { data, isFetching } = useQuery(
    projectSearchQueryOptions(
      {
        query: query || undefined,
        offset: page * perPage,
        limit: perPage,
      },
      {
        placeholderData: (previousData) => previousData,
        initialData: initialData,
        initialDataUpdatedAt: 0,
      },
    ),
  );

  return (
    <>
      <ProjectSearch
        total={data?.totalHits}
        query={_query ?? ""}
        setQuery={setQuery}
        clearSearchParams={() => {
          setQuery(null);
          setPage(null);
          setPerPage(null);
        }}
      />
      <ProjectList projects={data?.data} loading={isFetching} />
    </>
  );
};

export default Inner;
