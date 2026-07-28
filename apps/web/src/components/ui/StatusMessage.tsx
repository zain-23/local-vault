import {
  CircleAlert,
  CircleCheck,
  Info,
  type LucideIcon,
  TriangleAlert,
} from "lucide-react";
import type * as React from "react";

import { cn } from "#/lib/utils.ts";

type Tone = "success" | "error" | "info" | "warn";

const toneClass: Record<Tone, string> = {
  success: "border-success/25 bg-success/10 text-success",
  error: "border-destructive/25 bg-destructive/10 text-destructive",
  info: "border-primary/25 bg-primary/10 text-foreground",
  warn: "border-warning/40 bg-warning/10 text-foreground",
};

const iconClass: Record<Tone, string> = {
  success: "size-4",
  error: "size-4",
  info: "mt-0.5 size-3.5 shrink-0 text-primary",
  warn: "mt-0.5 size-3.5 shrink-0 text-warning",
};

type StatusMessageProps = {
  tone: Tone;
  icon: LucideIcon;
  /** Short status line (success / error). */
  title?: string;
  /** Multi-line body (info / warn banners). */
  children?: React.ReactNode;
  className?: string;
};

function StatusMessageRoot({
  tone,
  icon: Icon,
  title,
  children,
  className,
}: StatusMessageProps) {
  const isBanner = tone === "info" || tone === "warn";

  return (
    <div
      className={cn(
        "flex gap-2.5 rounded-lg border px-3.5 py-3 text-[12.5px]",
        isBanner ? "items-start" : "items-center",
        toneClass[tone],
        className,
      )}
    >
      <Icon className={cn("shrink-0", iconClass[tone])} />
      {title ? <span className="font-semibold">{title}</span> : null}
      {children ? <span className="leading-relaxed">{children}</span> : null}
    </div>
  );
}

type TitleProps = {
  title: string;
  className?: string;
};

type BannerProps = {
  children: React.ReactNode;
  className?: string;
};

function SuccessMessage({ title, className }: TitleProps) {
  return (
    <StatusMessageRoot
      tone="success"
      icon={CircleCheck}
      title={title}
      className={className}
    />
  );
}

function ErrorMessage({ title, className }: TitleProps) {
  return (
    <StatusMessageRoot
      tone="error"
      icon={CircleAlert}
      title={title}
      className={className}
    />
  );
}

function InfoBanner({ children, className }: BannerProps) {
  return (
    <StatusMessageRoot tone="info" icon={Info} className={className}>
      {children}
    </StatusMessageRoot>
  );
}

function WarnBanner({ children, className }: BannerProps) {
  return (
    <StatusMessageRoot tone="warn" icon={TriangleAlert} className={className}>
      {children}
    </StatusMessageRoot>
  );
}

const StatusMessage = Object.assign(StatusMessageRoot, {
  Success: SuccessMessage,
  Error: ErrorMessage,
  InfoBanner,
  WarnBanner,
});

export { StatusMessage, SuccessMessage, ErrorMessage };
