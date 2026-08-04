import { createFileRoute } from "@tanstack/react-router";

import { DashboardPage } from "#/features/dashboard/components";

export const Route = createFileRoute("/_app/dashboard")({
	staticData: { breadcrumb: "Dashboard" },
	component: DashboardPage,
});
