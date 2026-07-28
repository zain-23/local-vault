import { zodResolver } from "@hookform/resolvers/zod";
import { Link } from "@tanstack/react-router";
import { useForm } from "react-hook-form";

import {
  Button,
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
  Input,
} from "#/components/ui/index.ts";
import { type LoginValues, loginSchema } from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { GoogleButton } from "./GoogleButton.tsx";
import { PasswordField } from "./PasswordField.tsx";
import { useLogin } from "../hooks/useLogin.ts";

function LoginForm() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const { mutate, isPending } = useLogin();

  // UI only: validation runs, but there is no backend call yet.
  const onSubmit = handleSubmit((values) => {
    mutate(values);
  });

  return (
    <>
      <AuthHeading
        title="Welcome back"
        subtitle="Sign in to your LocalVault workspaces."
      />

      <GoogleButton />

      <FieldSeparator className="my-5">OR</FieldSeparator>

      <form onSubmit={onSubmit} noValidate>
        <FieldGroup className="gap-4">
          <Field data-invalid={!!errors.email}>
            <FieldLabel htmlFor="email">Email</FieldLabel>
            <Input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="ada@company.com"
              aria-invalid={!!errors.email}
              {...register("email")}
            />
            <FieldError errors={[errors.email]} />
          </Field>

          <Field data-invalid={!!errors.password}>
            <div className="flex items-center justify-between">
              <FieldLabel htmlFor="password">Password</FieldLabel>
              <Link
                to="/auth/forgot-password"
                className="text-[11.5px] text-muted-foreground hover:text-foreground"
              >
                Forgot password?
              </Link>
            </div>
            <PasswordField
              id="password"
              autoComplete="current-password"
              placeholder="••••••••"
              aria-invalid={!!errors.password}
              {...register("password")}
            />
            <FieldError errors={[errors.password]} />
          </Field>

          <Button type="submit" className="w-full" isLoading={isPending}>
            Log in
          </Button>
        </FieldGroup>
      </form>

      <AuthFooter>
        Need an account?{" "}
        <Link
          to="/auth/signup"
          className="font-medium text-foreground hover:underline"
        >
          Sign up
        </Link>
      </AuthFooter>
    </>
  );
}

export { LoginForm };
