import { parseAsString, parseAsInteger, createLoader } from "nuqs/server";

export const projectSearchParams = {
  query: parseAsString.withDefault(""),
  page: parseAsInteger.withDefault(0),
  perPage: parseAsInteger.withDefault(10),
};
export const loadProjectSearchParams = createLoader(projectSearchParams);
