import { cache } from "react";
import { client } from "./client";
import type { CreateProjectSchema } from "./models/project/create";
import type { Project, ProjectIdentifier, ProjectMember } from "./types";

const apiService = {
  project: {
    create: async (values: CreateProjectSchema): Promise<Project> => {
      const { data, error } = await client.POST("/projects", {
        body: {
          ...values,
        },
      });

      if (error) {
        throw Error(error.detail);
      }

      return data;
    },
    identifier: cache(
      async (identifier: ProjectIdentifier): Promise<Project | null> => {
        const { data, error } = await client.GET("/projects/{id|slug}", {
          params: {
            path: { "id|slug": identifier },
          },
        });

        if (error) {
          switch (error.status) {
            case 404:
              return null;
            default:
              throw new Error(error.detail);
          }
        }

        return data;
      },
    ),
    members: cache(
      async (
        identifier: ProjectIdentifier,
      ): Promise<Array<ProjectMember> | null> => {
        const { data, error } = await client.GET(
          "/projects/{id|slug}/members",
          {
            params: {
              path: { "id|slug": identifier },
            },
          },
        );

        if (error) {
          switch (error.status) {
            case 404:
              return null;
            default:
              throw new Error(error.detail);
          }
        }

        return data;
      },
    ),
  },
};

export default apiService;
