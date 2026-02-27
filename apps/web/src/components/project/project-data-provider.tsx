"use client";
import { Project, ProjectMember } from "@/lib/api/types";
import { createContext, useContext, type ReactNode } from "react";

type ProjectDataContextType = {
  project: Project | null;
  members: ProjectMember[] | null;
};

const ProjectDataContext = createContext<ProjectDataContextType | null>(null);

type ProjectDataProviderProps = {
  project: Project | null;
  members: ProjectMember[] | null;
  children: ReactNode;
};

export const ProjectDataProvider = ({
  project,
  members,
  children,
}: ProjectDataProviderProps) => {
  return (
    <ProjectDataContext.Provider value={{ project, members }}>
      {children}
    </ProjectDataContext.Provider>
  );
};

export const useProjectData = () => {
  const context = useContext(ProjectDataContext);
  if (!context) {
    throw new Error("useProjectData must be used within a ProjectDataProvider");
  }
  return context;
};

export default ProjectDataProvider;
