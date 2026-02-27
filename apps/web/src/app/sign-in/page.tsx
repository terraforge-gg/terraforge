"use client";
import { useForm } from "@tanstack/react-form";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { InputGroup, InputGroupInput } from "@/components/ui/input-group";
import { signIn } from "@/lib/auth-client";
import { signInSchema } from "@/lib/api/models/user/auth";
import { Icons } from "@/components/icons";
import { useRouter } from "next/navigation";
import Link from "@/components/link";

const SignInPage = () => {
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
        router.push("/");
      } else {
        toast.error(error.message ?? "Something went wrong.");
      }
    },
  });
  return (
    <div className="flex justify-center items-center">
      <Card className="w-96">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl">Sign up</CardTitle>
          <CardDescription>
            New to terraforge?{" "}
            <Link href="/sign-up" className="text-primary">
              Sign up for an account
            </Link>
            .
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4">
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
              selector={(state) => [
                state.isSubmitting,
                state.isSubmitSuccessful,
              ]}
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
              <span className="px-2 text-muted-foreground">
                Or continue with
              </span>
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
        </CardContent>
      </Card>
    </div>
  );
};

export default SignInPage;
