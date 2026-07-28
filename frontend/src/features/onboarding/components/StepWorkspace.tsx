import { ArrowRight } from "lucide-react";
import { useState } from "react";

import { Button, Field, FieldLabel, Input } from "#/components/ui/index.ts";
import { StepHeading } from "./OnboardingLayout.tsx";

// Turn a workspace name into the URL slug shown live under the input.
function slugify(name: string) {
  return name
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

// Step 1 — name the workspace. Slug preview is display-only (no availability check).
function StepWorkspace({ onContinue }: { onContinue: () => void }) {
  const [name, setName] = useState("Kodexo Labs");
  const slug = slugify(name) || "your-workspace";

  return (
    <>
      <StepHeading
        title="Create your workspace"
        subtitle="Workspaces hold your team, vaults, and audit history."
      />

      <Field>
        <FieldLabel htmlFor="workspace-name">Workspace name</FieldLabel>
        <Input
          id="workspace-name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Acme Inc"
        />
      </Field>

      <p className="mt-2.5 text-xs text-muted-foreground">
        URL:{" "}
        <span className="font-mono">
          localvault.app/<span className="text-foreground">{slug}</span>
        </span>
      </p>

      <Button className="mt-6 w-full" onClick={onContinue}>
        Continue
        <ArrowRight className="size-3.5" />
      </Button>
    </>
  );
}

export { StepWorkspace };
