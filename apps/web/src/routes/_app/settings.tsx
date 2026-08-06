import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { SettingsPage } from "#/features/settings/components";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/settings")({
	head: () => seo(PAGE_META["/_app/settings"]),
	staticData: { breadcrumb: "Settings" },
	component: SettingsPage,
});
