import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { ForgotPasswordForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/forgot-password")({
	head: () => seo(PAGE_META["/auth/forgot-password"]),
	component: ForgotPasswordForm,
});
