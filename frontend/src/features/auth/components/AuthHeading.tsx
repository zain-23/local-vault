import type * as React from "react";

import { LVLogo } from "#/components/shared/LVLogo.tsx";

// Header block above each auth form: a single top mark — the brand logo by
// default, or a contextual icon badge on confirmation-style screens — then the
// title and optional subtitle. Keeps every screen visually consistent.
function AuthHeading({
  title,
  subtitle,
}: {
  title: React.ReactNode;
  subtitle: React.ReactNode;
}) {
  return (
    <div className="mb-6 flex flex-col items-center text-center gap-y-2">
      <LVLogo size={26} withWord={false} />
      <h1 className="text-4xl font-semibold tracking-[-0.02em] text-balance">
        {title}
      </h1>
      <p className="text-sm leading-relaxed text-pretty text-muted-foreground">
        {subtitle}
      </p>
    </div>
  );
}

export { AuthHeading };
