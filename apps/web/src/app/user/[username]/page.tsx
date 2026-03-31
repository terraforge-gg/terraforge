import ProjectList from "@/components/project/project-list";
import UserHeader from "@/components/user/user-header";
import api from "@/lib/api/api";

const UserPage = async ({
  params,
}: {
  params: Promise<{ username: string }>;
}) => {
  const { username } = await params;
  const projects = await api.user.list(username);
  const userId = projects[0].userId;

  return (
    <>
      {/* temp */}
      <UserHeader username={userId} />
      <ProjectList projects={projects} />
    </>
  );
};

export default UserPage;
