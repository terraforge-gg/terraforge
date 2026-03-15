"use client";
import React from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { InputGroup, InputGroupInput } from "@/components/ui/input-group";
import { Icons } from "@/components/icons";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { signIn } from "@/lib/auth/client";
import { signInSchema } from "@/lib/auth/models/sign-in";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";

const SignInDialog = () => {
  const router = useRouter();
  const form = useForm({
    defaultValues: {
      username: "",
      password: "",
    },
    validators: {
      onSubmit: signInSchema,
    },
    onSubmit: async ({ value }) => {
      const { error } = await signIn.username({
        username: value.username,
        password: value.password,
      });
      form.reset();

      if (!error) {
        router.refresh();
      } else {
        toast.error(error.message ?? "Something went wrong.");
      }
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Sign in</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-sm">
        <DialogHeader>
          <DialogTitle>Sign in</DialogTitle>
        </DialogHeader>
        <form
          id="sign-in-form"
          className="space-y-4"
          onSubmit={(e) => {
            e.preventDefault();
            form.handleSubmit();
          }}
        >
          <FieldGroup>
            <form.Field
              name="username"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Username</FieldLabel>
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
              name="password"
              children={(field) => {
                const isInvalid =
                  field.state.meta.isTouched && !field.state.meta.isValid;
                return (
                  <Field>
                    <FieldLabel htmlFor={field.name}>Password</FieldLabel>
                    <InputGroup>
                      <InputGroupInput
                        id={field.name}
                        name={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        aria-invalid={isInvalid}
                        autoComplete="off"
                        type="password"
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
          <form.Subscribe
            selector={(state) => [state.isSubmitting, state.isSubmitSuccessful]}
            children={([isSubmitting]) => (
              <Field>
                <Button type="submit" disabled={isSubmitting}>
                  Sign In
                </Button>
              </Field>
            )}
          />
        </form>
        <div className="flex justify-center">
          <div className="flex justify-center text-xs uppercase">
            <span className="px-2 text-muted-foreground">Or continue with</span>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-6">
          <Button
            variant="outline"
            onClick={async () => {
              await signIn.social({
                provider: "discord",
              });
            }}
          >
            <Icons.discord className="mr-2 h-4 w-4" />
            Discord
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default SignInDialog;
