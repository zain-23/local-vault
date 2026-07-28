import type * as React from "react";

import { LVLogo } from "#/components/shared/LVLogo.tsx";
import { Stepper } from "./Stepper.tsx";

// Shared shell for every onboarding step: grid backdrop, centered logo, and a
// card that renders the Stepper once so the steps only supply their own body.
function OnboardingLayout({
  step,
  children,
}: {
  step: number;
  children: React.ReactNode;
}) {
  return (
    <div className="relative flex min-h-svh items-center justify-center overflow-hidden bg-background px-6 py-10 text-foreground">
      {/* grid motif, masked to a soft radial fade behind the card */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 opacity-70 bg-[linear-gradient(var(--color-border-strong)_1px,transparent_1px),linear-gradient(90deg,var(--color-border-strong)_1px,transparent_1px)] bg-size-[24px_24px]"
        style={{
          maskImage:
            "radial-gradient(ellipse at center, black 0%, transparent 70%)",
          WebkitMaskImage:
            "radial-gradient(ellipse at center, black 0%, transparent 70%)",
        }}
      />

      <div className="relative w-full max-w-xl">
        <div className="mb-6 flex justify-center">
          <LVLogo size={40} />
        </div>

        <div className="rounded-xl border border-border bg-card p-7">
          <Stepper current={step} />
          {children}
        </div>
      </div>
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
