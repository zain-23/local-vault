import { Link } from "@tanstack/react-router";
import { RefreshCw } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";

// Shared "check your inbox" confirmation for both the password-reset and
// magic-link flows — the copy switches on `variant`.
function CheckEmailPanel({
  variant,
  email = "your email",
}: {
  variant: "reset" | "magic";
  email?: string;
}) {
  const subtitle =
    variant === "magic" ? (
      <>
        We sent a magic sign-in link to{" "}
        <strong className="text-foreground">{email}</strong>. Open it on this
        device to continue.
      </>
    ) : (
      <>
        We sent a password reset link to{" "}
        <strong className="text-foreground">{email}</strong>. The link expires
        in 15 minutes.
      </>
    );

  return (
    <>
      <AuthHeading title="Check your inbox" subtitle={subtitle} />

      <Button type="button" size="lg" className="w-full">
        <RefreshCw />
        Resend email
      </Button>

      <AuthFooter>
        Wrong address?{" "}
        <Link
          to={
            variant === "magic" ? "/auth/magic-link" : "/auth/forgot-password"
          }
          className="font-medium text-foreground hover:underline"
        >
          Use a different one
        </Link>
      </AuthFooter>
    </>
  );
}

export { CheckEmailPanel };
