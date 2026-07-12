import { Link } from "@tanstack/react-router";
import { RefreshCw } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";

// Post-signup email verification prompt.
function VerifyEmailPanel({ email = "your email" }: { email?: string }) {
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
