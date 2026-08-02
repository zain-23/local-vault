import { createFileRoute, redirect } from "@tanstack/react-router";

import { meQuery } from "#/features/auth/api";
import { workspacesQuery } from "#/features/onboarding/api";
import { OnboardingWizard } from "#/features/onboarding/components/index.ts";

export const Route = createFileRoute("/onboarding")({
	ssr: false,
	beforeLoad: async ({ context }) => {
		const user = await context.queryClient.ensureQueryData(meQuery);
		if (!user) throw redirect({ to: "/auth/login" });
		if (user.onboarded) throw redirect({ to: "/" });
	},
	// Prefetch the caller's workspaces so Step 1 can prefill without a flash and
	// resume from an already-created workspace after a refresh.
	loader: ({ context }) => context.queryClient.ensureQueryData(workspacesQuery),
	component: OnboardingWizard,
});
