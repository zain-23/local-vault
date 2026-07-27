import type { ReactNode } from "react";

import { cn } from "#/lib/utils.ts";

type SettingsBlockProps = {
  label: string;
  description?: string;
  children: ReactNode;
  danger?: boolean;
  className?: string;
};

// Preference row: label + description on the left, controls on the right.
// Stacks on small screens so Cursor-density stays readable on mobile.
export function SettingsBlock({
  label,
  description,
  children,
  danger = false,
  className,
}: SettingsBlockProps) {
  return (
    <div
      className={cn(
        "grid grid-cols-1 gap-4 border-b border-border py-5 sm:grid-cols-[minmax(0,1fr)_minmax(200px,280px)] sm:gap-6",
        className,
      )}
    >
      <div className="min-w-0">
        <div
          className={cn(
            "text-sm font-medium",
            danger ? "text-destructive" : "text-foreground",
          )}
        >
          {label}
        </div>
        {description ? (
          <p className="mt-1 max-w-sm text-[13px] leading-relaxed text-muted-foreground">
            {description}
          </p>
        ) : null}
      </div>
      <div className="flex min-w-0 flex-col gap-2.5">{children}</div>
    </div>
  );
}
