import { Link } from "@tanstack/react-router";
import { RefreshCw } from "lucide-react";
import { useEffect } from "react";

import {
  Button,
  ErrorMessage,
  Spinner,
  SuccessMessage,
} from "#/components/ui/index.ts";
import { useVerifyEmail } from "../hooks/useVerifyEmail.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";

function VerifyEmailPanel({
  email = "your email",
  token,
}: {
  email?: string;
  token?: string;
}) {
  const { mutate, isPending, isSuccess, isError } = useVerifyEmail();

  // Arrived from the email link — verify the token once on mount.
  useEffect(() => {
    if (token) {
      mutate(token);
    }
  }, [token, mutate]);

  if (token) {
    return (
      <>
        <AuthHeading
          title="Verify your email"
          subtitle="Confirming your email address — this only takes a moment."
        />

        {isPending && (
          <div className="flex flex-col items-center gap-3 py-4 text-muted-foreground">
            <Spinner size="lg" />
            <p className="text-sm">Verifying your link…</p>
          </div>
        )}

        {isSuccess && (
          <SuccessMessage
            title="Email verified"
          />
        )}

        {isError && (
          <ErrorMessage
            title="Verification failed"
          />
        )}
      </>
    );
  }

  // No token: the post-signup "check your inbox" prompt.
  return (
    <>
      <AuthHeading
        title="Verify your email"
        subtitle={
          <>
            We sent a verification link to{" "}
            <strong className="text-foreground">{email}</strong>. Click it to
            activate your account.
          </>
        }
      />

      <Button type="button" size="lg" className="w-full">
        <RefreshCw />
        Resend verification
      </Button>

      <AuthFooter>
        Already verified?{" "}
        <Link
          to="/auth/login"
          className="font-medium text-foreground hover:underline"
        >
          Log in
        </Link>
      </AuthFooter>
    </>
  );
}

export { VerifyEmailPanel };
