"use client";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { InputGroup, InputGroupInput } from "@/components/ui/input-group";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import api from "@/lib/api/api";
import { createProjectReleaseSchema } from "@/lib/api/models/project/create-release";
import { loaderVersionsQueryOptions } from "@/lib/api/query-options/loader-version";
import { projectReleasePresignedPutUrlQueryOptions } from "@/lib/api/query-options/project-release";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChangeEvent, useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";

type CreateProjectReleaseDialogProps = {
  projectId: string;
  projectSlug: string;
};

const CreateProjectReleaseDialog = ({
  projectId,
  projectSlug,
}: CreateProjectReleaseDialogProps) => {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();
  const defaultValues: z.input<typeof createProjectReleaseSchema> = {
    name: "",
    versionNumber: "",
    loaderVersionId: "",
    changelog: undefined,
    fileUrl: "",
    dependencies: [],
  };
  const [file, setFile] = useState<File | undefined>(undefined);
  const form = useForm({
    defaultValues: defaultValues,
    validators: {
      onSubmit: createProjectReleaseSchema,
    },
    onSubmit: ({ value }) => {
      createProjectVersion({
        projectId,
        values: {
          ...value,
        },
      });
    },
  });

  const {
    mutate: createProjectVersion,
    isPending: createProjectVersionPending,
  } = useMutation({
    mutationFn: api.project.createRelease,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["project-releases", { projectSlug: projectSlug }],
      });
      form.reset();
      toast.success("Project release created!");
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  const { data: loaderVersions, isPending } = useQuery(
    loaderVersionsQueryOptions(),
  );

  const { data: presignedPutUrl } = useQuery(
    projectReleasePresignedPutUrlQueryOptions(
      {
        projectId,
        fileSize: file?.size ?? 0,
      },
      {
        enabled: !!file,
      },
    ),
  );

  const { mutate: uploadProjectVersionFile, isPending: isUploadPending } =
    useMutation<void, Error, { presignedPutUrl: string; file: File }>({
      mutationKey: ["upload-project-version", projectId],
      mutationFn: async ({ presignedPutUrl: url, file: f }) => {
        await fetch(url, {
          method: "PUT",
          body: f,
          headers: {
            "content-type": "application/octet-stream",
            "content-length": `${file?.size ?? 0}`,
          },
        });
      },
      onSuccess: () => {
        setFile(undefined);
        if (presignedPutUrl) {
          const newIconUrl = new URL(presignedPutUrl);
          form.setFieldValue(
            "fileUrl",
            newIconUrl.origin + newIconUrl.pathname,
          );
          toast.success("File uploaded");
        }
      },
      onError: () => {
        toast.error("Something went wrong");
        setFile(undefined);
      },
    });

  useEffect(() => {
    if (!presignedPutUrl || !file) return;
    uploadProjectVersionFile({ presignedPutUrl, file });
  }, [presignedPutUrl]);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newFile = e.target.files?.[0];

    if (!newFile) return;

    setFile(newFile);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">CREATE RELEASE</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Create Release</DialogTitle>
        </DialogHeader>
        <form
          id="create-project-version-form"
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <Input id="picture" type="file" onChange={handleFileChange} />
          <FieldGroup>
            <form.Field
              name="name"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Name</FieldLabel>
                    <InputGroup>
                      <InputGroupInput
                        id={field.name}
                        name={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        autoComplete="off"
                        className="pl-0!"
                      />
                    </InputGroup>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
          </FieldGroup>
          <FieldGroup>
            <form.Field
              name="versionNumber"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Version</FieldLabel>
                    <InputGroup>
                      <InputGroupInput
                        id={field.name}
                        name={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        autoComplete="off"
                        className="pl-0!"
                      />
                    </InputGroup>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
          </FieldGroup>
          <FieldGroup>
            <form.Field
              name="changelog"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Changelog</FieldLabel>
                    <Textarea
                      id={field.name}
                      name={field.name}
                      value={field.state.value ?? undefined}
                      onBlur={field.handleBlur}
                      onChange={(e) => field.handleChange(e.target.value)}
                      aria-invalid={isInvalid}
                      autoComplete="off"
                    />
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
          </FieldGroup>
          <FieldGroup>
            <form.Field
              name="loaderVersionId"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field className="max-w-56">
                    <FieldLabel htmlFor={field.name}>Loader Version</FieldLabel>
                    <Select onValueChange={(id) => field.handleChange(id)}>
                      <SelectTrigger
                        className="w-full"
                        id={field.name}
                        disabled={isPending}
                      >
                        <SelectValue
                          placeholder={
                            isPending ? "Loading..." : "Select a loader version"
                          }
                        />
                      </SelectTrigger>
                      <SelectContent position="popper">
                        <SelectGroup>
                          {loaderVersions?.map((x) => (
                            <SelectItem key={x.id} value={x.id}>
                              {x.gameVersion} - ({x.buildType}) - (
                              {x.versionLabel})
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                    {isInvalid && (
                      <FieldError errors={field.state.meta.errors} />
                    )}
                  </Field>
                );
              }}
            />
          </FieldGroup>
          <Field className="flex w-full justify-end">
            <Button className="w-full" disabled={createProjectVersionPending}>
              Submit
            </Button>
          </Field>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default CreateProjectReleaseDialog;
