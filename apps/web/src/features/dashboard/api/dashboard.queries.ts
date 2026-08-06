import { queryOptions } from "@tanstack/react-query";
import type { ActivityParams } from "#/features/dashboard/api";
import { DASHBOARD_KEYS, dashboardService } from "#/features/dashboard/api";

export function dashboardSummaryQuery(workspaceId: string) {
	return queryOptions({
		queryKey: DASHBOARD_KEYS.summary(workspaceId),
		queryFn: async () => {
			const res = await dashboardService.summary(workspaceId);
			return res.data;
		},
		staleTime: 30_000,
	});
}

export function dashboardActivityQuery(
	workspaceId: string,
	params: ActivityParams = {},
) {
	return queryOptions({
		queryKey: DASHBOARD_KEYS.activity(workspaceId, params),
		queryFn: async () => {
			const res = await dashboardService.activity(workspaceId, params);
			return res.data;
		},
		staleTime: 20_000,
	});
}
