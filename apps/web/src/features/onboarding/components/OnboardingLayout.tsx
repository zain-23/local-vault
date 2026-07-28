import type * as React from "react";

import { LVLogo } from "#/components/shared/LVLogo.tsx";
import { Stepper } from "./Stepper.tsx";

const FOOTER_LINKS = ["Privacy", "Terms", "Security"] as const;

function OnboardingLayout({
  step,
  children,
}: {
  step: number;
  children: React.ReactNode;
}) {
  return (
    <div className="relative flex min-h-svh flex-col overflow-hidden bg-background text-foreground">
      {/* Grid backdrop, matched to the device + auth screens so the flows feel
          of a piece. Masked to a soft glow anchored top-left. */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 opacity-60 bg-[linear-gradient(var(--color-border)_1px,transparent_1px),linear-gradient(90deg,var(--color-border)_1px,transparent_1px)] bg-size-[24px_24px]"
        style={{
          maskImage:
            "radial-gradient(ellipse at 25% 15%, black 0%, transparent 72%)",
          WebkitMaskImage:
            "radial-gradient(ellipse at 25% 15%, black 0%, transparent 72%)",
        }}
      />

      <header className="relative flex items-center justify-between px-6 py-5 sm:px-10">
        <LVLogo size={35} />
        <button
          type="button"
          className="text-[13px] text-muted-foreground transition-colors hover:text-foreground"
        >
          Need help?
        </button>
      </header>

      <main className="relative flex flex-1 items-center justify-center px-6 py-10">
        <div className="w-full max-w-xl">
          <div className="rounded-xl border border-border bg-card p-7">
            <Stepper current={step} />
            {children}
          </div>
        </div>
      </main>

      <footer className="relative flex flex-wrap items-center justify-center gap-x-4 gap-y-1 px-6 py-6 text-[11.5px] text-muted-foreground">
        <span className="font-mono">© 2026 LocalVault</span>
        {FOOTER_LINKS.map((label) => (
          <button
            key={label}
            type="button"
            className="transition-colors hover:text-foreground"
          >
            {label}
          </button>
        ))}
      </footer>
    </div>
  );
}

// Title + subtitle block shared by every step body.
function StepHeading({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <div className="mb-5.5">
      <h2 className="text-[22px] font-semibold tracking-[-0.02em]">{title}</h2>
      <p className="mt-1.5 text-[13px] text-muted-foreground">{subtitle}</p>
    </div>
  );
}

export { OnboardingLayout, StepHeading };
