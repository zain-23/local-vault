import { createFileRoute } from "@tanstack/react-router";

import { ResetPasswordForm } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/reset-password")({
	component: ResetPasswordForm,
});
