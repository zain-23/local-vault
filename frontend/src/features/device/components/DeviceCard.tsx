import type { LucideIcon } from "lucide-react";
import type * as React from "react";

// The card every device screen shares: a tinted icon badge, a title, an optional
// subtitle, then the screen's body. Keeps the submit and approval steps visually
// identical apart from their contents.
function DeviceCard({
  icon: Icon,
  title,
  subtitle,
  header,
  children,
}: {
  icon?: LucideIcon;
  title?: React.ReactNode;
  subtitle?: React.ReactNode;
  // Optional custom header (used by the approval screen's side-by-side layout).
  header?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="rounded-2xl border border-border bg-card/60 p-7 shadow-xs backdrop-blur-sm sm:p-8">
      {header ?? (
        <div className="mb-6 flex flex-col items-center gap-3 text-center">
          {Icon && (
            <div className="flex size-12 items-center justify-center rounded-2xl border border-border bg-muted text-primary">
              <Icon className="size-5.5" />
            </div>
          )}
          {title && (
            <h1 className="text-2xl font-semibold tracking-[-0.02em] text-balance">
              {title}
            </h1>
          )}
          {subtitle && (
            <p className="text-sm leading-relaxed text-pretty text-muted-foreground">
              {subtitle}
            </p>
          )}
        </div>
      )}
      {children}
    </div>
  );
}

export { DeviceCard };
