"use client";
import { useProjectData } from "@/components/project/project-data-provider";
import { notFound } from "next/navigation";

const ModPage = () => {
  const { project: mod, members } = useProjectData();

  if (!mod || !members) {
    notFound();
  }

  if (!mod.description) return null;

  return (
    <section className="flex flex-col gap-4">
      <h2 className="font-mono text-xs">DESCRIPTION</h2>
      <p className="text-muted-foreground">{mod.description}</p>
    </section>
  );
};

export default ModPage;
