"use client";
import ProjectList from "@/components/project/project-list";
import ProjectSearch from "@/components/project/project-search";
import { Button } from "@/components/ui/button";
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
      <div className="w-80 sm:w-150">
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
      </div>
      <ProjectList projects={data?.data} loading={isFetching} />
      <div className="flex w-full justify-end gap-2">
        {page > 0 && (
          <Button variant="outline" onClick={() => setPage((prev) => prev - 1)}>
            Previous
          </Button>
        )}
        {perPage < (data?.totalHits || 0) / perPage && (
          <Button variant="outline" onClick={() => setPage((prev) => prev + 1)}>
            Next
          </Button>
        )}
      </div>
    </>
  );
};

export default Inner;
