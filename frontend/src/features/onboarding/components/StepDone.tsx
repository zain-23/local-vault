import { Button } from "#/components/ui";
import { useCompleteOnboarding } from "../hooks/index.ts";
import { NEXT_STEPS } from "../utils/constant.ts";
import { StepHeading } from "./OnboardingLayout.tsx";

// Onboarding is a CLI-first product, so the payoff of "done" is the first real
// commands to run — not a recap of setup. These are UI-only (mock) prompts.
function StepDone() {
  const { isPending, mutate } = useCompleteOnboarding();

  return (
    <>
      <StepHeading
        title="You're all set"
        subtitle="Your workspace is live. Here's where to start from your terminal."
      />

      {/* lighter follow-up commands */}
      <div className="flex flex-col gap-2">
        {NEXT_STEPS.map((item) => (
          <div
            key={item.label}
            className="flex items-center gap-2.5 rounded-md border border-border px-3 py-2.5"
          >
            <item.icon className="size-3.5 shrink-0 text-muted-foreground" />
            <span className="flex-1 text-[13px] font-medium">{item.label}</span>
            <code className="font-mono text-[11.5px] text-muted-foreground">
              {item.command}
            </code>
          </div>
        ))}
      </div>

      <Button
        className="mt-6 w-full"
        isLoading={isPending}
        onClick={() => mutate()}
      >
        Go to dashboard
      </Button>
    </>
  );
}

export { StepDone };
