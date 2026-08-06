import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { ResetSuccessPanel } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/reset-success")({
	head: () => seo(PAGE_META["/auth/reset-success"]),
	component: ResetSuccessPanel,
});
