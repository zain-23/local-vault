import type * as React from "react";

// Centered muted line under a form (e.g. "Need an account? Sign up").
function AuthFooter({ children }: { children: React.ReactNode }) {
  return (
    <p className="mt-5.5 text-center text-[13px] text-muted-foreground">
      {children}
    </p>
  );
}

export { AuthFooter };
