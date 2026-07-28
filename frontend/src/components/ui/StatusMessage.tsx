import type * as React from "react";
import { CircleAlert, CircleCheck, type LucideIcon } from "lucide-react";

import { cn } from "#/lib/utils.ts";

type Tone = "success" | "error";

function StatusMessage({
  tone,
  icon: Icon,
  title,
}: {
  tone: Tone;
  icon: LucideIcon;
  title: string;
}) {
  return (
    <div
      className={cn(
        "flex items-center gap-2.5 rounded-lg border px-3.5 py-3 text-[12.5px]",
        tone === "success"
          ? "border-success/25 bg-success/10 text-success"
          : "border-destructive/25 bg-destructive/10 text-destructive",
      )}
    >
      <Icon className="size-4 shrink-0" />
      <span className="font-semibold">{title}</span>
    </div>
  );
}

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
