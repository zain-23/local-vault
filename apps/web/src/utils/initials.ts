export function initialsFromName(name: string): string {
	const parts = name.trim().split(/\s+/).filter(Boolean);
	if (parts.length === 0) return "?";
	const first = parts[0] ?? "";
	if (parts.length === 1) return first.slice(0, 2).toUpperCase();
	const second = parts[1] ?? "";
	return `${first.charAt(0)}${second.charAt(0)}`.toUpperCase();
}
