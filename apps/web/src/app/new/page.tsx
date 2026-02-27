import CreateProjectForm from "@/components/project/create-project-form";
import getServerSession from "@/lib/auth";
import { notFound } from "next/navigation";

const Page = async () => {
  const session = await getServerSession();

  if (!session) {
    notFound();
  }

  return (
    <div className="flex justify-center">
      <CreateProjectForm />
    </div>
  );
};

export default Page;
