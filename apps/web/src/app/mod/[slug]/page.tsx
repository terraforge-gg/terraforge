"use client";
import { useProjectData } from "@/components/project/project-data-provider";
import { cn } from "@/lib/utils";
import { notFound } from "next/navigation";

const ModPage = () => {
  const { project: mod, members } = useProjectData();

  if (!mod || !members) {
    notFound();
  }

  return (
    <section className="flex flex-col gap-4">
      <p
        className={cn("text-muted-foreground", mod.description ? "" : "italic")}
      >
        {mod.description ?? "no description"}
      </p>
    </section>
  );
};

export default ModPage;
