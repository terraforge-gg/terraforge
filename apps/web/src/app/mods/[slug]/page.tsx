"use client";
import { useProjectData } from "@/components/project/project-data-provider";
import { Card, CardContent } from "@/components/ui/card";
import { notFound } from "next/navigation";

const ModPage = () => {
  const { project: mod, members } = useProjectData();

  if (!mod || !members) {
    notFound();
  }

  return (
    <Card>
      <CardContent className="px-4">
        {mod.description ?? (
          <span className="italic text-muted-foreground">no description</span>
        )}
      </CardContent>
    </Card>
  );
};

export default ModPage;
