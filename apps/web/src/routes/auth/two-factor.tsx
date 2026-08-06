import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { TwoFactorForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/two-factor")({
	head: () => seo(PAGE_META["/auth/two-factor"]),
	component: TwoFactorForm,
});
