import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { AuditPage } from "#/features/audit/components";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/audit")({
	head: () => seo(PAGE_META["/_app/audit"]),
	staticData: { breadcrumb: "Audit log" },
	component: AuditPage,
});
