import { type ApiClient, api } from "#/services/api";

import type {
	ActivityParams,
	DashboardActivity,
	DashboardSummary,
} from "./dashboard.types.ts";

class DashboardService {
	constructor(private readonly client: ApiClient = api) {}

	private base(workspaceId: string) {
		return `/workspaces/${encodeURIComponent(workspaceId)}/dashboard`;
	}

	summary(workspaceId: string) {
		return this.client.get<DashboardSummary>(
			`${this.base(workspaceId)}/summary`,
		);
	}

	activity(workspaceId: string, params: ActivityParams = {}) {
		return this.client.get<DashboardActivity>(
			`${this.base(workspaceId)}/activity`,
			{ params: { range: params.range } },
		);
	}
}

export const dashboardService = new DashboardService();
export { DashboardService };
