import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { authGuard } from "#/features/auth/utils/authGuard";
import { workspacesQuery } from "#/features/onboarding/api";
import { OnboardingWizard } from "#/features/onboarding/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/onboarding")({
	head: () => seo(PAGE_META["/onboarding"]),
	...authGuard((user) => {
		if (!user) return { to: "/auth/login" };
		if (user.onboarded) return { to: "/dashboard" };
	}),
	// Prefetch the caller's workspaces so Step 1 can prefill without a flash and
	// resume from an already-created workspace after a refresh.
	loader: ({ context }) => context.queryClient.ensureQueryData(workspacesQuery),
	component: OnboardingWizard,
});
