import { createFileRoute } from "@tanstack/react-router";

import { AuditPage } from "#/features/audit/components";

export const Route = createFileRoute("/_app/audit")({
	component: AuditPage,
});
