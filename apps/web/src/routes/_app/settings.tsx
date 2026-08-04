import { createFileRoute } from "@tanstack/react-router";

import { SettingsPage } from "#/features/settings/components";

export const Route = createFileRoute("/_app/settings")({
	staticData: { breadcrumb: "Settings" },
	component: SettingsPage,
});
