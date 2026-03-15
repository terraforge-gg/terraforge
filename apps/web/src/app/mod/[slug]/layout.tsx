import React from "react";
import { notFound } from "next/navigation";
import { Project, PROJECT_MEMBER_ROLE } from "@/lib/api/types";
import getServerSession from "@/lib/auth/session";
import ProjectDataProvider from "@/components/project/project-data-provider";
import api from "@/lib/api/api";
import ProjectHeader from "@/components/project/project-header";
import ProjectMembers from "@/components/project/project-members";
import ProjectReleases from "@/components/project/project-releases";

const ModLayout = async ({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ slug: string }>;
}) => {
  const { slug } = await params;
  const session = await getServerSession();
  const mod = await api.project.identifier(slug);
  const members = await api.project.members(slug);
  const releases = await api.project.releases(slug);

  if (!mod || !members) {
    notFound();
  }

  const role = members?.find((x) => x.userId === session?.user.id)?.role;
  const canViewSettings =
    role === PROJECT_MEMBER_ROLE.OWNER || role === PROJECT_MEMBER_ROLE.ADMIN;

  return (
    <ProjectDataProvider project={mod} members={members}>
      <div className="flex min-h-screen gap-10">
        <div className="flex min-w-0 flex-1 flex-col gap-6">
          <ProjectHeader project={mod} latestRelease={releases[0]} />
          <section className="flex flex-col gap-4">
            <a
              href="#summary"
              className="w-16 font-mono text-xs transition ease-in-out hover:text-chart-1"
            >
              SUMMARY
            </a>
            <p className="text-muted-foreground">{mod.summary}</p>
          </section>
          {children}
        </div>
        <aside className="w-80 shrink-0">
          <div className="sticky top-14 flex flex-col gap-4 pt-10">
            <ProjectMembers members={members} />
            <ProjectReleases releases={releases} />
          </div>
        </aside>
      </div>
    </ProjectDataProvider>
  );
};

export default ModLayout;
