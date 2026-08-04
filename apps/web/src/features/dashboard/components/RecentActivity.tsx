import { Link } from "@tanstack/react-router";

import { Skeleton } from "#/components/ui";
import type { AuditEvent } from "#/features/audit/api";
import { AuditEventRow } from "#/features/audit/components/AuditEventRow.tsx";

type RecentActivityProps = {
	events: AuditEvent[] | undefined;
	isLoading: boolean;
};

export function RecentActivity({ events, isLoading }: RecentActivityProps) {
	if (isLoading || !events) {
		return (
			<div className="space-y-3">
				{(["r1", "r2", "r3"] as const).map((key) => (
					<div key={key} className="flex gap-3 py-2">
						<Skeleton className="h-4 w-14" />
						<Skeleton className="size-7 rounded-md" />
						<div className="flex-1 space-y-2">
							<Skeleton className="h-4 w-2/3" />
							<Skeleton className="h-3 w-1/4" />
						</div>
					</div>
				))}
			</div>
		);
	}

	if (events.length === 0) {
		return (
			<p className="py-6 text-sm text-muted-foreground">
				No recent events.{" "}
				<Link
					to="/audit"
					className="font-medium text-foreground underline-offset-4 hover:underline"
				>
					Open audit log
				</Link>
			</p>
		);
	}

	return (
		<div>
			{events.map((event) => (
				<AuditEventRow key={event.id} event={event} />
			))}
			<div className="border-t border-border pt-3">
				<Link
					to="/audit"
					className="text-sm text-muted-foreground underline-offset-4 hover:text-foreground hover:underline"
				>
					View audit log
				</Link>
			</div>
		</div>
	);
}
