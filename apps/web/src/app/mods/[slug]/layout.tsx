import React from "react";
import {
  Item,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemHeader,
  ItemMedia,
  ItemTitle,
} from "@/components/ui/item";
import UserAvatar from "@/components/user/user-avatar";
import ProjectHeader from "@/components/project/project-header";
import ProjectNavbar from "@/components/project/project-navbar";
import apiService from "@/lib/api/service";
import { notFound } from "next/navigation";
import Link from "next/link";
import ProjectDataProvider from "@/components/project/project-data-provider";
import getServerSession from "@/lib/auth";
import { PROJECT_MEMBER_ROLE } from "@/lib/api/types";

const ModLayout = async ({
  children,
  params,
}: {
  children: React.ReactNode;
  params: Promise<{ slug: string }>;
}) => {
  const { slug } = await params;
  const session = await getServerSession();
  const mod = await apiService.project.identifier(slug);
  const members = await apiService.project.members(slug);

  if (!mod || !members) {
    notFound();
  }

  const role = members?.find((x) => x.userId === session?.user.id)?.role;
  const canViewSettings =
    role === PROJECT_MEMBER_ROLE.OWNER || role === PROJECT_MEMBER_ROLE.ADMIN;

  return (
    <ProjectDataProvider project={mod} members={members}>
      <ProjectHeader
        name={mod.name}
        summary={mod.summary}
        iconUrl={mod.iconUrl}
        downloads={mod.downloads}
      />
      <ProjectNavbar slug={mod.slug} showSettings={canViewSettings} />
      <div className="flex flex-col gap-4 md:flex-row">
        <div className="w-full min-h-96">{children}</div>
        <div className="flex flex-col w-full md:w-auto gap-4 md:ml-auto">
          <ItemGroup>
            {members.map((user) => (
              <React.Fragment key={user.username}>
                <Item variant="muted" className="w-full md:w-72">
                  <ItemHeader className="text-lg font-semibold">
                    Members
                  </ItemHeader>
                  <ItemMedia>
                    <UserAvatar
                      avatar={user.image}
                      fallback={user.username.charAt(0)}
                      className="w-10 h-10"
                    />
                  </ItemMedia>
                  <ItemContent className="gap-1">
                    <ItemTitle>
                      <Link
                        href="/sign-in"
                        className="hover:underline underline-offset-4 hover:decoration-primary"
                      >
                        {user.username}
                      </Link>
                    </ItemTitle>
                    <ItemDescription>
                      {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                    </ItemDescription>
                  </ItemContent>
                </Item>
              </React.Fragment>
            ))}
          </ItemGroup>
        </div>
      </div>
    </ProjectDataProvider>
  );
};

export default ModLayout;
