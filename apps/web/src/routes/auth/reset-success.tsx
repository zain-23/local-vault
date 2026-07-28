import { createFileRoute } from "@tanstack/react-router";

import { ResetSuccessPanel } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/reset-success")({
	component: ResetSuccessPanel,
});
