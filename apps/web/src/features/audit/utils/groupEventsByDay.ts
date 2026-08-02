import type { AuditEvent } from "#/features/audit/api";

export type AuditDayGroup = {
	/** YYYY-MM-DD local date key */
	key: string;
	label: string;
	events: AuditEvent[];
};

export function groupEventsByDay(
	events: AuditEvent[],
	now = new Date(),
): AuditDayGroup[] {
	const groups = new Map<string, AuditEvent[]>();

	for (const event of events) {
		const d = new Date(event.created_at);
		if (Number.isNaN(d.getTime())) continue;
		const key = localDateKey(d);
		const list = groups.get(key);
		if (list) list.push(event);
		else groups.set(key, [event]);
	}

	return [...groups.entries()]
		.sort(([a], [b]) => (a < b ? 1 : a > b ? -1 : 0))
		.map(([key, dayEvents]) => ({
			key,
			label: dayLabel(key, now),
			events: dayEvents,
		}));
}

function localDateKey(d: Date): string {
	const y = d.getFullYear();
	const m = String(d.getMonth() + 1).padStart(2, "0");
	const day = String(d.getDate()).padStart(2, "0");
	return `${y}-${m}-${day}`;
}

function dayLabel(key: string, now: Date): string {
	const today = localDateKey(now);
	const yesterdayDate = new Date(now);
	yesterdayDate.setDate(yesterdayDate.getDate() - 1);
	const yesterday = localDateKey(yesterdayDate);

	if (key === today) return "Today";
	if (key === yesterday) return "Yesterday";

	const parts = key.split("-");
	const y = Number(parts[0]);
	const m = Number(parts[1]);
	const d = Number(parts[2]);
	const date = new Date(y, m - 1, d);
	return date.toLocaleDateString(undefined, {
		weekday: "short",
		month: "short",
		day: "numeric",
		year: "numeric",
	});
}
