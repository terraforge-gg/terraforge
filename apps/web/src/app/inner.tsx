"use client";
import ProjectList from "@/components/project/project-list";
import ProjectSearch from "@/components/project/project-search";
import { Button } from "@/components/ui/button";
import { projectSearchQueryOptions } from "@/lib/api/query-options/project";
import { ProjectSearch as ProjectSearchResult } from "@/lib/api/types";
import { useQuery } from "@tanstack/react-query";
import { parseAsInteger, parseAsString, useQueryState, useQueryStates } from "nuqs";
import { useDebounce } from "use-debounce";

type SearchInnerProps = {
  initialData?: ProjectSearchResult;
};

const Inner = ({ initialData }: SearchInnerProps) => {
  const [params, setParams] = useQueryStates({
    query: parseAsString,
    page: parseAsInteger.withDefault(0),
    perPage: parseAsInteger.withDefault(10),
  })

  const [_query, setQuery] = useQueryState("query");
  const [query] = useDebounce(params.query, 500);

  const { data, isFetching } = useQuery(
    projectSearchQueryOptions(
      {
        query: query || undefined,
        offset: params.page * params.perPage,
        limit: params.perPage,
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
      <div className="w-80 sm:w-150">
        <ProjectSearch
          query={_query ?? ""}
          setQuery={setQuery}
          clearSearchParams={() => {
            setParams({
              query: null,
              page: null,
              perPage: null,
            })
          }}
        />
      </div>
      <ProjectList projects={data?.data} loading={isFetching} />
      <div className="flex w-full justify-end gap-2">
        {params.page > 0 && (
          <Button variant="outline" onClick={() => setParams((prev) => ({
            page: prev.page - 1
          }))}>
            Previous
          </Button>
        )}
        {params.perPage < (data?.totalHits || 0) / params.perPage && (
          <Button variant="outline" onClick={() => setParams((prev) => ({
            page: prev.page + 1
          }))}>
            Next
          </Button>
        )}
      </div>
    </>
  );
};

export default Inner;
