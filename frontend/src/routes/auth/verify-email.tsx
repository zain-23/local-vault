import { createFileRoute } from "@tanstack/react-router";

import { VerifyEmailPanel } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/verify-email")({
	component: VerifyEmailPanel,
});
