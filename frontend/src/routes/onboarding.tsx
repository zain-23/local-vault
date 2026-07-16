import { createFileRoute, redirect } from "@tanstack/react-router";

import { OnboardingWizard } from "#/features/onboarding/components/index.ts";
import { meQuery } from "#/features/auth/api";

export const Route = createFileRoute("/onboarding")({
  beforeLoad: async ({ context }) => {
    const user = await context.queryClient.ensureQueryData(meQuery);
    if (!user) throw redirect({ to: "/auth/login" });
    if (user.onboarded) throw redirect({ to: "/" });
  },
  component: OnboardingWizard,
});
