import React from "react";
import { notFound } from "next/navigation";
import { PROJECT_MEMBER_ROLE } from "@/lib/api/types";
import getServerSession from "@/lib/auth/session";
import ProjectDataProvider from "@/components/project/project-data-provider";
import api from "@/lib/api/api";
import ProjectHeader from "@/components/project/project-header";
import ProjectMembers from "@/components/project/project-members";
import ProjectReleases from "@/components/project/project-releases";
import { Separator } from "@/components/ui/separator";

const ModLayout = async ({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ slug: string }>;
}) => {
  const { slug } = await params;
  const session = await getServerSession();

  let mod = undefined;
  let members = undefined;
  let releases = undefined;

  try {
    mod = await api.project.identifier(slug);
    members = await api.project.members(slug);
    releases = await api.project.releases({ projectSlug: slug });
  } catch (error) {
    console.error(error);
    throw new Error("Failed to fetch mod info");
  }

  if (!mod || !members) {
    notFound();
  }

  const role = members?.find((x) => x.userId === session?.user.id)?.role;
  const canViewSettings =
    role === PROJECT_MEMBER_ROLE.OWNER || role === PROJECT_MEMBER_ROLE.ADMIN;
  const canCreateRelease =
    role === PROJECT_MEMBER_ROLE.OWNER || role === PROJECT_MEMBER_ROLE.ADMIN;

  return (
    <ProjectDataProvider project={mod} members={members}>
      <ProjectHeader project={mod} showSettings={canViewSettings} />
      <div className="flex min-h-screen gap-10">
        <div className="flex min-w-0 flex-1 flex-col gap-6">{children}</div>
        <aside className="w-80 shrink-0">
          <div className="sticky top-20 flex flex-col gap-8">
            <ProjectMembers members={members} />
            <Separator />
            <ProjectReleases
              projectId={mod.id}
              projectSlug={mod.slug}
              initialReleases={releases}
              showCreateRelease={canCreateRelease}
            />
          </div>
        </aside>
      </div>
    </ProjectDataProvider>
  );
};

export default ModLayout;
