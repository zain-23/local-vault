export function formatEventClock(iso: string): string {
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return "—";
	return d.toLocaleTimeString(undefined, {
		hour: "2-digit",
		minute: "2-digit",
		hour12: false,
	});
}

export function initialsFromName(name: string): string {
	const parts = name.trim().split(/\s+/).filter(Boolean);
	if (parts.length === 0) return "?";
	const first = parts[0] ?? "";
	if (parts.length === 1) return first.slice(0, 2).toUpperCase();
	const second = parts[1] ?? "";
	return `${first.charAt(0)}${second.charAt(0)}`.toUpperCase();
}

/** Stable hue from a string for avatar color (0–360). */
export function hueFromString(value: string): number {
	let hash = 0;
	for (let i = 0; i < value.length; i++) {
		hash = (hash * 31 + value.charCodeAt(i)) | 0;
	}
	return Math.abs(hash) % 360;
}
