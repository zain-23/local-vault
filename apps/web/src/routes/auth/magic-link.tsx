import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { MagicLinkForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/magic-link")({
	head: () => seo(PAGE_META["/auth/magic-link"]),
	component: MagicLinkForm,
});
