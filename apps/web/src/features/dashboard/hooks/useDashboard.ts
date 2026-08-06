import { useQuery } from "@tanstack/react-query";
import type { ActivityParams } from "#/features/dashboard/api";
import {
	dashboardActivityQuery,
	dashboardSummaryQuery,
} from "#/features/dashboard/api";
import { useWorkspaceStore } from "#/stores";

export function useDashboardSummary() {
	const workspaceId = useWorkspaceStore((s) => s.active?.id);
	return useQuery({
		...dashboardSummaryQuery(workspaceId ?? ""),
		enabled: !!workspaceId,
	});
}

export function useDashboardActivity(params: ActivityParams = {}) {
	const workspaceId = useWorkspaceStore((s) => s.active?.id);
	return useQuery({
		...dashboardActivityQuery(workspaceId ?? "", params),
		enabled: !!workspaceId,
	});
}
