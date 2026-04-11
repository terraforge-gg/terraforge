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
import {
  createProjectReleaseSchema,
  PROJECT_RELEASE_MAX_DEPENDENCIES,
} from "@/lib/api/models/project/create-release";
import {
  PROJECT_RELEASE_DEPENDENCY_TYPES,
  type ProjectReleaseDependencyType,
} from "@/lib/api/types";
import { loaderVersionsQueryOptions } from "@/lib/api/query-options/loader-version";
import { projectReleasePresignedPutUrlQueryOptions } from "@/lib/api/query-options/project-release";
import { useForm } from "@tanstack/react-form";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChangeEvent, useEffect, useState } from "react";
import { toast } from "sonner";
import z from "zod";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";

type CreateProjectReleaseDialogProps = {
  projectId: string;
  projectSlug: string;
};

const CreateProjectReleaseDialog = ({
  projectId,
  projectSlug,
}: CreateProjectReleaseDialogProps) => {
  const [open, setOpen] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [showAllVersions, setShowAllVersions] = useState(false);
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
  const [depProjectId, setDepProjectId] = useState("");
  const [depType, setDepType] = useState<ProjectReleaseDependencyType>(
    PROJECT_RELEASE_DEPENDENCY_TYPES[0],
  );
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
      setOpen(false);
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

  const { mutate: uploadProjectVersionFile } = useMutation<
    void,
    Error,
    { presignedPutUrl: string; file: File }
  >({
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
        form.setFieldValue("fileUrl", newIconUrl.origin + newIconUrl.pathname);
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
          <Input
            id="picture"
            type="file"
            accept=".tmod"
            onChange={handleFileChange}
          />
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

                let filteredLoaderVersions = showPreview
                  ? loaderVersions
                  : loaderVersions?.filter((x) => x.buildType === "stable");

                if (!showAllVersions && filteredLoaderVersions) {
                  filteredLoaderVersions = Object.values(
                    filteredLoaderVersions.reduce<
                      Record<string, (typeof filteredLoaderVersions)[number]>
                    >((acc, version) => {
                      if (
                        !acc[version.gameVersion] ||
                        new Date(version.releasedAt) >
                          new Date(acc[version.gameVersion].releasedAt)
                      ) {
                        acc[version.gameVersion] = version;
                      }
                      return acc;
                    }, {}),
                  );
                }

                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Loader Version</FieldLabel>
                    <div className="flex flex-col gap-2">
                      <Select onValueChange={(id) => field.handleChange(id)}>
                        <SelectTrigger
                          className="w-full"
                          id={field.name}
                          disabled={isPending}
                        >
                          <SelectValue
                            placeholder={
                              isPending
                                ? "Loading..."
                                : "Select a loader version"
                            }
                          />
                        </SelectTrigger>
                        <SelectContent position="popper">
                          <SelectGroup>
                            {filteredLoaderVersions?.map((x) => (
                              <SelectItem key={x.id} value={x.id}>
                                <div className="flex items-center gap-1.5">
                                  <span>{x.gameVersion}</span>
                                  <span>({x.versionLabel})</span>
                                  {x.buildType === "preview" && (
                                    <Badge>{x.buildType}</Badge>
                                  )}
                                </div>
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </SelectContent>
                      </Select>
                      <div className="flex items-center gap-3">
                        <div className="flex items-center gap-1.5">
                          <Checkbox
                            id="show-preview"
                            checked={showPreview}
                            onCheckedChange={(checked) =>
                              setShowPreview(checked === true)
                            }
                          />
                          <label
                            htmlFor="show-preview"
                            className="cursor-pointer text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                          >
                            Show preview
                          </label>
                        </div>
                        <div className="flex items-center gap-1.5">
                          <Checkbox
                            id="show-all-versions"
                            checked={showAllVersions}
                            onCheckedChange={(checked) =>
                              setShowAllVersions(checked === true)
                            }
                          />
                          <label
                            htmlFor="show-all-versions"
                            className="cursor-pointer text-sm leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                          >
                            Show all versions
                          </label>
                        </div>
                      </div>
                    </div>
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
              name="dependencies"
              children={(field) => {
                const deps = field.state.value ?? [];
                const canAdd = deps.length < PROJECT_RELEASE_MAX_DEPENDENCIES;
                return (
                  <Field>
                    <FieldLabel>Dependencies</FieldLabel>
                    <div className="flex gap-2">
                      <InputGroup className="flex-1">
                        <InputGroupInput
                          placeholder="Project ID"
                          value={depProjectId}
                          onChange={(e) => setDepProjectId(e.target.value)}
                          autoComplete="off"
                        />
                      </InputGroup>
                      <Select
                        value={depType}
                        onValueChange={(v) =>
                          setDepType(v as ProjectReleaseDependencyType)
                        }
                      >
                        <SelectTrigger className="w-28">
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent position="popper">
                          <SelectGroup>
                            {PROJECT_RELEASE_DEPENDENCY_TYPES.map((t) => (
                              <SelectItem key={t} value={t}>
                                {t}
                              </SelectItem>
                            ))}
                          </SelectGroup>
                        </SelectContent>
                      </Select>
                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        disabled={!canAdd || !depProjectId}
                        onClick={() => {
                          field.pushValue({
                            type: depType,
                            projectId: depProjectId,
                          });
                          setDepProjectId("");
                        }}
                      >
                        Add
                      </Button>
                    </div>
                    {deps.length > 0 && (
                      <div className="mt-1.5 flex flex-wrap gap-1.5">
                        {deps.map((dep, i) => (
                          <Badge
                            key={`${dep.projectId}-${i}`}
                            variant="secondary"
                            className="gap-1 pr-1"
                          >
                            <span className="max-w-32 truncate">
                              {dep.projectId}
                            </span>
                            <span className="opacity-60">{dep.type}</span>
                            <button
                              type="button"
                              className="ml-0.5 hover:text-destructive"
                              onClick={() => field.removeValue(i)}
                            >
                              ×
                            </button>
                          </Badge>
                        ))}
                      </div>
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
