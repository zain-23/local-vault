import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { DashboardPage } from "#/features/dashboard/components";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/dashboard")({
	head: () => seo(PAGE_META["/_app/dashboard"]),
	staticData: { breadcrumb: "Dashboard" },
	component: DashboardPage,
});
