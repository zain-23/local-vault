import type { SeriesPoint } from "#/features/dashboard/api";

export function maxSeriesTotal(series: SeriesPoint[]): number {
	let max = 0;
	for (const point of series) {
		if (point.total > max) max = point.total;
	}
	return max;
}

/** Short label for a YYYY-MM-DD UTC date (tooltip). */
export function formatSeriesDay(date: string): string {
	const [y, m, d] = date.split("-").map(Number);
	if (!y || !m || !d) return date;
	const dt = new Date(Date.UTC(y, m - 1, d));
	return dt.toLocaleDateString(undefined, {
		weekday: "short",
		month: "short",
		day: "numeric",
		year: "numeric",
		timeZone: "UTC",
	});
}

/** GitHub-style intensity 0–4 from count vs range max. */
export function activityLevel(total: number, max: number): 0 | 1 | 2 | 3 | 4 {
	if (total <= 0 || max <= 0) return 0;
	const ratio = total / max;
	if (ratio <= 0.25) return 1;
	if (ratio <= 0.5) return 2;
	if (ratio <= 0.75) return 3;
	return 4;
}

export const ACTIVITY_LEVEL_CLASS: Record<0 | 1 | 2 | 3 | 4, string> = {
	0: "bg-muted",
	1: "bg-primary/25",
	2: "bg-primary/45",
	3: "bg-primary/70",
	4: "bg-primary",
};

export type HeatCell = {
	date: string;
	total: number;
	inRange: boolean;
};

export type MonthLabel = {
	weekIndex: number;
	label: string;
};

const MONTH_SHORT = [
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
] as const;

/** Pad series into Sunday-start week columns for a contribution calendar. */
export function buildContributionWeeks(series: SeriesPoint[]): HeatCell[][] {
	if (series.length === 0) return [];

	const byDate = new Map(series.map((p) => [p.date, p.total]));
	const first = parseUtcDate(series[0].date);
	const last = parseUtcDate(series[series.length - 1].date);
	if (!first || !last) return [];

	const start = new Date(first);
	start.setUTCDate(start.getUTCDate() - start.getUTCDay());

	const end = new Date(last);
	end.setUTCDate(end.getUTCDate() + (6 - end.getUTCDay()));

	const cells: HeatCell[] = [];
	for (
		const cursor = new Date(start);
		cursor <= end;
		cursor.setUTCDate(cursor.getUTCDate() + 1)
	) {
		const key = formatUtcKey(cursor);
		const inRange = cursor >= first && cursor <= last;
		cells.push({
			date: key,
			inRange,
			total: inRange ? (byDate.get(key) ?? 0) : 0,
		});
	}

	const weeks: HeatCell[][] = [];
	for (let i = 0; i < cells.length; i += 7) {
		weeks.push(cells.slice(i, i + 7));
	}
	return weeks;
}

/** Month labels placed above the first week where that month appears. */
export function buildMonthLabels(weeks: HeatCell[][]): MonthLabel[] {
	const labels: MonthLabel[] = [];
	let prevMonth: number | null = null;

	weeks.forEach((week, weekIndex) => {
		const cell = week.find((c) => c.inRange) ?? week[0];
		if (!cell) return;
		const d = parseUtcDate(cell.date);
		if (!d) return;
		const month = d.getUTCMonth();
		if (prevMonth === null || month !== prevMonth) {
			labels.push({ weekIndex, label: MONTH_SHORT[month] });
			prevMonth = month;
		}
	});

	return labels;
}

function parseUtcDate(date: string): Date | null {
	const [y, m, d] = date.split("-").map(Number);
	if (!y || !m || !d) return null;
	return new Date(Date.UTC(y, m - 1, d));
}

function formatUtcKey(d: Date): string {
	const y = d.getUTCFullYear();
	const m = String(d.getUTCMonth() + 1).padStart(2, "0");
	const day = String(d.getUTCDate()).padStart(2, "0");
	return `${y}-${m}-${day}`;
}
