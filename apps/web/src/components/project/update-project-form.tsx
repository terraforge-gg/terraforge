"use client";
import { useForm, useStore } from "@tanstack/react-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type z from "zod";
import type { Project } from "@/lib/api/types";
import { getChangedFields } from "@/lib/utils";
import apiService from "@/lib/api/service";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
  InputGroupTextarea,
} from "@/components/ui/input-group";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { updateProjectSchema } from "@/lib/api/models/project/update";
import type { UpdateProjectParams } from "@/lib/api/models/project/update";
import { Button } from "../ui/button";
import { useRouter } from "next/navigation";
import { env } from "@/env";

type UpdateProjectFormProps = {
  project: Project;
};

const UpdateProjectForm = ({ project }: UpdateProjectFormProps) => {
  const queryClient = useQueryClient();
  const router = useRouter();
  const defaultValues: z.input<typeof updateProjectSchema> = {
    name: project.name,
    slug: project.slug,
    summary: project.summary ?? undefined,
    iconUrl: project.iconUrl ?? undefined,
  };

  const form = useForm({
    defaultValues,
    validators: {
      onSubmit: updateProjectSchema,
    },
    onSubmit: ({ value }) => {
      updateProject({
        projectIdentifier: project.id,
        values: getChangedFields(value, defaultValues),
      });
    },
  });

  const isDefaultValue = useStore(form.store, (state) => state.isDefaultValue);

  const { mutate: updateProject, isPending: isUpdatePending } = useMutation({
    mutationFn: (variables: UpdateProjectParams) => {
      return apiService.project.update({
        projectIdentifier: variables.projectIdentifier,
        values: variables.values,
      });
    },
    onSuccess: async () => {
      const newSlug = form.getFieldValue("slug");
      form.reset();
      if (newSlug && newSlug != defaultValues.slug) {
        router.push(`/mods/${newSlug}/settings`);
      } else {
        queryClient.invalidateQueries({
          queryKey: ["project", { identifier: project.slug }],
        });

        router.refresh();
      }
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  return (
    <div className="flex flex-col gap-4">
      <Card>
        <CardHeader>
          <CardTitle>Project Information</CardTitle>
        </CardHeader>
        <CardContent>
          <form
            id="update-project-form"
            className="flex flex-col gap-4"
            onSubmit={(e) => {
              e.preventDefault();
              form.handleSubmit();
            }}
          >
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
                name="summary"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field data-invalid={isInvalid}>
                      <FieldLabel htmlFor={field.name}>Summary</FieldLabel>
                      <InputGroup>
                        <InputGroupTextarea
                          id={field.name}
                          name={field.name}
                          value={field.state.value}
                          onBlur={field.handleBlur}
                          onChange={(e) => field.handleChange(e.target.value)}
                          aria-invalid={isInvalid}
                          autoComplete="off"
                        />
                        <InputGroupAddon align="block-end">
                          <InputGroupText className="ml-auto">{`${field.state.value?.length ?? 0}/120`}</InputGroupText>
                        </InputGroupAddon>
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
                name="slug"
                children={(field) => {
                  const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid;
                  return (
                    <Field>
                      <FieldLabel htmlFor={field.name}>Slug</FieldLabel>
                      <InputGroup>
                        <InputGroupAddon>
                          <InputGroupText>{`${env.NEXT_PUBLIC_APP_URL}/mods/`}</InputGroupText>
                        </InputGroupAddon>
                        <InputGroupInput
                          className="pl-0!"
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
            <Field className="w-32">
              <Button
                type="submit"
                disabled={isDefaultValue || isUpdatePending}
              >
                Save Changes
              </Button>
            </Field>
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default UpdateProjectForm;
