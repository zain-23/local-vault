import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { AuthHeading } from "./AuthHeading.tsx";

// Terminal confirmation after a successful password change.
function ResetSuccessPanel() {
  return (
    <>
      <AuthHeading
        title="Password updated"
        subtitle="Your password has been changed. Sign in with your new password to continue."
      />

      <Button asChild className="w-full">
        <Link to="/auth/login">
          Continue to login
          <ArrowRight />
        </Link>
      </Button>
    </>
  );
}

export { ResetSuccessPanel };
