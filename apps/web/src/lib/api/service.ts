import { client } from "./client";
import type { CreateProjectSchema } from "./models/project/create";
import type { Project } from "./types";

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
  },
};

export default apiService;
