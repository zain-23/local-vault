import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { SignupForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/signup")({
	head: () => seo(PAGE_META["/auth/signup"]),
	component: SignupForm,
});
