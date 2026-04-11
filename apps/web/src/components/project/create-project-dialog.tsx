import { PlusIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  createProjectSchema,
  PROJECT_SUMMARY_MAX_LENGTH,
} from "@/lib/api/models/project/create";
import z from "zod";
import { useRouter } from "next/navigation";
import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
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
import { env } from "@/env";
import api from "@/lib/api/api";
import { useState } from "react";

const CreateProjectDialog = () => {
  const router = useRouter();
  const [open, setOpen] = useState(false);

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
    mutationFn: api.project.create,
    onSuccess: (data) => {
      router.push(`/mod/${data.slug}`);
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <PlusIcon />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Create Project</DialogTitle>
        </DialogHeader>
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
                        disabled={isPending}
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
                        <InputGroupText>
                          {env.NEXT_PUBLIC_APP_URL}/mod/
                        </InputGroupText>
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
                        disabled={isPending}
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
                    <InputGroup>
                      <InputGroupTextarea
                        id={field.name}
                        name={field.name}
                        value={field.state.value ?? undefined}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        autoComplete="off"
                        disabled={isPending}
                      />
                      <InputGroupAddon align="block-end">
                        <InputGroupText className="ml-auto">{`${field.state.value?.length ?? 0}/${PROJECT_SUMMARY_MAX_LENGTH}`}</InputGroupText>
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
          <Field className="flex w-full justify-end">
            <Button className="w-full" disabled={isPending}>
              Submit
            </Button>
          </Field>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default CreateProjectDialog;
