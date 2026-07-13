import type * as React from "react";
import { CircleAlert, CircleCheck, type LucideIcon } from "lucide-react";

import { cn } from "#/lib/utils.ts";

type Tone = "success" | "error";

// Centered status block for confirmation-style screens: a tinted icon badge, a
// title, an optional description, and an optional action area (buttons/links)
// passed as children.
function StatusMessage({
  tone,
  icon: Icon,
  title,
  description,
  className,
  children,
}: {
  tone: Tone;
  icon: LucideIcon;
  title: string;
  description?: React.ReactNode;
  className?: string;
  children?: React.ReactNode;
}) {
  return (
    <div
      className={cn(
        "flex flex-col items-center gap-4 rounded-2xl border border-border bg-card/40 px-6 py-8 text-center",
        className,
      )}
    >
      <div
        className={cn(
          "flex size-14 items-center justify-center rounded-full",
          tone === "success"
            ? "bg-success/10 text-success"
            : "bg-destructive/10 text-destructive",
        )}
      >
        <Icon className="size-7" />
      </div>

      <div className="space-y-1.5">
        <p className="text-lg font-semibold text-foreground">{title}</p>
        {description && (
          <p className="text-sm leading-relaxed text-muted-foreground">
            {description}
          </p>
        )}
      </div>

      {children && <div className="w-full pt-1">{children}</div>}
    </div>
  );
}

// Intent-named wrappers so call sites read clearly and pick the right icon/tone.
type MessageProps = Omit<
  React.ComponentProps<typeof StatusMessage>,
  "tone" | "icon"
>;

function SuccessMessage(props: MessageProps) {
  return <StatusMessage tone="success" icon={CircleCheck} {...props} />;
}

function ErrorMessage(props: MessageProps) {
  return <StatusMessage tone="error" icon={CircleAlert} {...props} />;
}

export { StatusMessage, SuccessMessage, ErrorMessage };
