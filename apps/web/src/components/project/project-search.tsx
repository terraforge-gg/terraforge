"use client";
import { SearchIcon, XIcon } from "lucide-react";
import {
  InputGroup,
  InputGroupInput,
  InputGroupAddon,
  InputGroupButton,
} from "@/components/ui/input-group";
import { Options } from "nuqs";

type ProjectSearchProps = {
  query?: string;
  setQuery: (
    value: string | ((old: string | null) => string | null) | null,
    options?: Options | undefined,
  ) => Promise<URLSearchParams>;
  clearSearchParams: () => void;
};

const ProjectSearch = ({
  query,
  setQuery,
  clearSearchParams,
}: ProjectSearchProps) => {
  return (
    <InputGroup className="h-12">
      <InputGroupInput
        placeholder="Search mods..."
        value={query ?? ""}
        onChange={(e) => setQuery(e.target.value || null)}
      />
      <InputGroupAddon>
        <SearchIcon />
      </InputGroupAddon>
      {query && (
        <InputGroupAddon align="inline-end">
          <InputGroupButton
            size="icon-sm"
            className="ml-auto"
            onClick={clearSearchParams}
          >
            <XIcon />
            <span className="sr-only">Clear</span>
          </InputGroupButton>
        </InputGroupAddon>
      )}
    </InputGroup>
  );
};

export default ProjectSearch;
