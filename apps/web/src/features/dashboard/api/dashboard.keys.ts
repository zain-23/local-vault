import type { ActivityParams } from "./dashboard.types.ts";

export const DASHBOARD_KEYS = {
	all: ["dashboard"] as const,
	workspace: (workspaceId: string) =>
		[...DASHBOARD_KEYS.all, workspaceId] as const,
	summary: (workspaceId: string) =>
		[...DASHBOARD_KEYS.workspace(workspaceId), "summary"] as const,
	activity: (workspaceId: string, params: ActivityParams = {}) =>
		[...DASHBOARD_KEYS.workspace(workspaceId), "activity", params] as const,
};
