import { createFileRoute } from "@tanstack/react-router";

import { ForgotPasswordForm } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/forgot-password")({
	component: ForgotPasswordForm,
});
