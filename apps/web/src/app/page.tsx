import api from "@/lib/api/api";
import { loadProjectSearchParams } from "@/lib/api/searchParams";
import { SearchParams } from "nuqs/server";
import Inner from "./inner";

type PageProps = {
  searchParams: Promise<SearchParams>;
};

const Home = async ({ searchParams }: PageProps) => {
  const { query, page, perPage } = await loadProjectSearchParams(searchParams);
  let result = undefined;
  try {
    result = await api.project.search({
      query,
      offset: page * perPage,
    });
  } catch (error) {
    console.error("Failed to search for projects.", error);
  }

  return (
    <div className="flex min-h-screen flex-col items-center gap-10">
      <h1 className="text-4xl">Find your favourite mods</h1>
      <Inner initialData={result} />
    </div>
  );
};

export default Home;
