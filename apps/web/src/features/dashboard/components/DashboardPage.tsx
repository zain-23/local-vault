import { Spinner } from "#/components/ui";
import {
	useDashboardActivity,
	useDashboardSummary,
} from "#/features/dashboard/hooks";
import { useWorkspaceStore } from "#/stores";

import { ActivityChart } from "./ActivityChart.tsx";
import { RecentActivity } from "./RecentActivity.tsx";
import { SummaryStats } from "./SummaryStats.tsx";

export function DashboardPage() {
	const workspace = useWorkspaceStore((s) => s.active);

	const summaryQuery = useDashboardSummary();
	const activityQuery = useDashboardActivity({ range: "1y" });

	const hasError = summaryQuery.isError || activityQuery.isError;
	const errorMessage =
		(summaryQuery.error instanceof Error && summaryQuery.error.message) ||
		(activityQuery.error instanceof Error && activityQuery.error.message) ||
		"Could not load dashboard.";

	const totalEvents =
		activityQuery.data?.series.reduce((sum, p) => sum + p.total, 0) ?? 0;

	return (
		<div className="flex flex-col gap-4 p-6">
			<div>
				<h1 className="text-xl font-semibold tracking-tight">Dashboard</h1>
				<p className="text-sm text-muted-foreground">
					Overview of{" "}
					{workspace?.name ? (
						<span className="text-foreground">{workspace.name}</span>
					) : (
						"your workspace"
					)}
					.
				</p>
			</div>

			{hasError ? (
				<div
					role="alert"
					className="rounded-lg border border-destructive/40 bg-destructive/5 px-4 py-3 text-sm text-destructive"
				>
					{errorMessage}
				</div>
			) : null}

			<SummaryStats
				summary={summaryQuery.data}
				isLoading={summaryQuery.isLoading}
			/>

			<section className="rounded-lg border border-border bg-background p-4">
				<div className="mb-4 flex items-center justify-between gap-3">
					<div>
						<h2 className="text-sm font-medium">Activity</h2>
						{!activityQuery.isLoading ? (
							<p className="mt-0.5 text-xs text-muted-foreground">
								{totalEvents} event{totalEvents === 1 ? "" : "s"} in the last
								year
							</p>
						) : null}
					</div>
					{activityQuery.isFetching && !activityQuery.isLoading ? (
						<Spinner className="size-3.5 text-muted-foreground" />
					) : null}
				</div>
				<ActivityChart
					series={activityQuery.data?.series}
					isLoading={activityQuery.isLoading}
				/>
			</section>

			<section className="rounded-lg border border-border bg-background p-4">
				<h2 className="mb-4 text-sm font-medium">Recent</h2>
				<RecentActivity
					events={activityQuery.data?.recent}
					isLoading={activityQuery.isLoading}
				/>
			</section>
		</div>
	);
}
