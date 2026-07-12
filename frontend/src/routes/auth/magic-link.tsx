import { createFileRoute } from "@tanstack/react-router";

import { MagicLinkForm } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/magic-link")({
	component: MagicLinkForm,
});
