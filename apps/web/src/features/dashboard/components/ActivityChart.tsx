import { useId, useMemo } from "react";

import { Skeleton } from "#/components/ui";
import type { SeriesPoint } from "#/features/dashboard/api";
import {
	ACTIVITY_LEVEL_CLASS,
	activityLevel,
	buildContributionWeeks,
	buildMonthLabels,
	formatSeriesDay,
	maxSeriesTotal,
} from "#/features/dashboard/utils";
import { cn } from "#/lib/utils.ts";

type ActivityChartProps = {
	series: SeriesPoint[] | undefined;
	isLoading: boolean;
};

const WEEKDAY_ROWS = [
	{ id: "sun", label: "" },
	{ id: "mon", label: "Mon" },
	{ id: "tue", label: "" },
	{ id: "wed", label: "Wed" },
	{ id: "thu", label: "" },
	{ id: "fri", label: "Fri" },
	{ id: "sat", label: "" },
] as const;

export function ActivityChart({ series, isLoading }: ActivityChartProps) {
	const labelId = useId();
	const max = useMemo(() => maxSeriesTotal(series ?? []), [series]);
	const weeks = useMemo(() => buildContributionWeeks(series ?? []), [series]);
	const months = useMemo(() => buildMonthLabels(weeks), [weeks]);
	const monthByWeek = useMemo(() => {
		const map = new Map<number, string>();
		for (const m of months) map.set(m.weekIndex, m.label);
		return map;
	}, [months]);

	if (isLoading || !series) {
		return (
			<div className="w-full space-y-2">
				<Skeleton className="h-3 w-full" />
				<Skeleton className="h-22 w-full" />
			</div>
		);
	}

	return (
		<div className="w-full">
			<div className="overflow-x-auto">
				<div
					role="img"
					aria-labelledby={labelId}
					className="flex min-w-160 gap-1.5 sm:min-w-0"
				>
					<span id={labelId} className="sr-only">
						Contribution calendar of workspace activity over the past year
					</span>

					{/* Weekday gutter — aligned under month row */}
					<div className="flex w-7 shrink-0 flex-col">
						<span className="mb-1 h-3" aria-hidden />
						<div className="flex flex-1 flex-col justify-between gap-px">
							{WEEKDAY_ROWS.map((row) => (
								<span
									key={row.id}
									className="flex flex-1 items-center text-[10px] leading-none text-muted-foreground"
								>
									{row.label}
								</span>
							))}
						</div>
					</div>

					{/* Weeks stretch to fill the box */}
					<div className="flex min-w-0 flex-1 gap-1.5">
						{weeks.map((week, weekIndex) => {
							const weekKey = week[0]?.date ?? `week-${weekIndex}`;
							const month = monthByWeek.get(weekIndex);
							return (
								<div key={weekKey} className="flex min-w-0 flex-1 flex-col">
									<span className="mb-1 h-3 text-[10px] leading-none text-muted-foreground">
										{month ?? ""}
									</span>
									<div className="flex flex-1 flex-col gap-1.5">
										{week.map((cell) => {
											if (!cell.inRange) {
												return (
													<span
														key={cell.date}
														className={cn(
															"aspect-square w-full min-w-0 rounded-[2px]",
															"bg-transparent",
														)}
														aria-hidden
													/>
												);
											}
											const level = activityLevel(cell.total, max);
											return (
												<span
													key={cell.date}
													title={`${formatSeriesDay(cell.date)}: ${cell.total} event${cell.total === 1 ? "" : "s"}`}
													className={cn(
														"aspect-square w-full min-w-0 rounded-[2px]",
														"ring-1 ring-border/50",
														ACTIVITY_LEVEL_CLASS[level],
													)}
												/>
											);
										})}
									</div>
								</div>
							);
						})}
					</div>
				</div>
			</div>
		</div>
	);
}
