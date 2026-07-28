import { ArrowRight } from "lucide-react";
import { Button, Field, FieldLabel, Input } from "#/components/ui/index.ts";
import { StepHeading } from "./OnboardingLayout.tsx";

function StepWorkspace({ onContinue }: { onContinue: () => void }) {
  return (
    <>
      <StepHeading
        title="Create your workspace"
        subtitle="Workspaces hold your team, vaults, and audit history."
      />

      <Field>
        <FieldLabel htmlFor="workspace-name">Workspace name</FieldLabel>
        <Input id="workspace-name" placeholder="Acme Inc" />
      </Field>

      <Button className="mt-6 w-full" onClick={onContinue}>
        Continue
        <ArrowRight className="size-3.5" />
      </Button>
    </>
  );
}

export { StepWorkspace };
