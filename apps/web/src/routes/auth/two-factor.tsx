import { createFileRoute } from "@tanstack/react-router";

import { TwoFactorForm } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/two-factor")({
	component: TwoFactorForm,
});
