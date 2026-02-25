import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useNavigate } from "@tanstack/react-router";
import type z from "zod";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
} from "@/components/ui/input-group";
import { Textarea } from "@/components/ui/textarea";
import { createProjectSchema } from "@/lib/api/models/project/create";
import apiService from "@/lib/api/service";
import { env } from "@/env/client";
import { Button } from "../ui/button";

const CreateProjectForm = () => {
  const navigate = useNavigate();

  const defaultValues: z.input<typeof createProjectSchema> = {
    type: "mod",
    name: "",
    slug: "",
    summary: undefined,
    organisationId: undefined,
  };

  const form = useForm({
    defaultValues: defaultValues,
    validators: {
      onSubmit: createProjectSchema,
    },
    onSubmit: ({ value }) => {
      createProject({
        ...value,
      });
    },
  });

  const { mutate: createProject, isPending } = useMutation({
    mutationFn: apiService.project.create,
    onSuccess: (data) => {
      form.reset();
      navigate({
        to: "/mod/$slug",
        params: { slug: data.slug },
      });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  return (
    <Card className="w-full max-w-150">
      <CardHeader>
        <CardTitle>Create Project</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          id="create-project-form"
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <FieldGroup className="flex flex-row">
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
              name="slug"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Slug</FieldLabel>
                    <InputGroup>
                      <InputGroupAddon>
                        <InputGroupText>{env.VITE_APP_URL}/mod/</InputGroupText>
                      </InputGroupAddon>
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
              name="summary"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Summary</FieldLabel>
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
          <Field className="w-full flex justify-end">
            <Button className="w-full" disabled={isPending}>
              Submit
            </Button>
          </Field>
        </form>
      </CardContent>
    </Card>
  );
};

export default CreateProjectForm;
