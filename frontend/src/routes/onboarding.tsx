import { createFileRoute } from "@tanstack/react-router";

import { OnboardingWizard } from "#/features/onboarding/components/index.ts";

export const Route = createFileRoute("/onboarding")({
	component: OnboardingWizard,
});
