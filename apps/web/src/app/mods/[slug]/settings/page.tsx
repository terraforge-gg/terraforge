"use client";
import { useProjectData } from "@/components/project/project-data-provider";
import apiService from "@/lib/api/service";
import { useMutation } from "@tanstack/react-query";
import { notFound, useRouter } from "next/navigation";
import { toast } from "sonner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import UpdateProjectForm from "@/components/project/update-project-form";
import { useSession } from "@/lib/auth-client";
import { Spinner } from "@/components/ui/spinner";

const ModSettingsPage = () => {
  const { project: mod, members } = useProjectData();
  const router = useRouter();
  const { data: session } = useSession();
  const { mutate: deleteProject, isPending: isDeletePending } = useMutation({
    mutationFn: apiService.project.delete,
    onSuccess: () => {
      toast.success("Project deleted");
      router.push("/");
    },
    onError: (error) => {
      toast.error(error.message ?? "Something went wrong.");
    },
  });

  if (!mod || !members) {
    notFound();
  }

  const role = members?.find((x) => x.userId === session?.user.id)?.role;
  const canDelete = role === "owner" || role === "admin";

  return (
    <div className="flex flex-col gap-4">
      <UpdateProjectForm project={mod} />
      {canDelete && (
        <Card>
          <CardHeader>
            <CardTitle>Danger</CardTitle>
          </CardHeader>
          <CardContent>
            <AlertDialog>
              <AlertDialogTrigger asChild>
                <Button variant="destructive">Delete</Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
                  <AlertDialogDescription>
                    This action cannot be undone.
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel disabled={isDeletePending}>
                    Cancel
                  </AlertDialogCancel>
                  <AlertDialogAction
                    onClick={() => deleteProject(mod.id)}
                    disabled={isDeletePending}
                  >
                    {isDeletePending && <Spinner />}
                    Confirm
                  </AlertDialogAction>
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default ModSettingsPage;
