import ProjectList from "@/components/project/project-list";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { SearchProjectParams } from "@/lib/api/models/project/search";
import apiService from "@/lib/api/service";
import { SearchIcon, XIcon } from "lucide-react";

export default async function HomePage() {
  const search: SearchProjectParams = {};

  let projects;

  try {
    projects = await apiService.project.search();
  } catch {}

  return (
    <div className="flex flex-col items-center gap-10 min-h-screen">
      <h1 className="text-4xl">Find your favourite mods</h1>
      <InputGroup className="w-150 h-12">
        <InputGroupInput
          placeholder={projects?.totalHits ? "Search mods..." : "No mods :("}
          defaultValue={search?.query ?? ""}
        />
        <InputGroupAddon>
          <SearchIcon />
        </InputGroupAddon>
        {search?.query && (
          <InputGroupAddon align="inline-end">
            <InputGroupButton size="icon-sm" className="ml-auto">
              <XIcon />
              <span className="sr-only">Clear</span>
            </InputGroupButton>
          </InputGroupAddon>
        )}
      </InputGroup>
      <ProjectList projects={projects?.data ?? []} />
    </div>
  );
}
