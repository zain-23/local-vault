export type ActionTone = "accent" | "success" | "warning" | "danger" | "muted";

const TONE_BY_VERB: Record<string, ActionTone> = {
	invited: "accent",
	invite: "accent",
	joined: "success",
	join: "success",
	created: "success",
	create: "success",
	pushed: "accent",
	push: "accent",
	pulled: "muted",
	pull: "muted",
	rotated: "warning",
	rotate: "warning",
	revoked: "danger",
	revoke: "danger",
	removed: "danger",
	remove: "danger",
	deleted: "danger",
	delete: "danger",
	updated: "muted",
	update: "muted",
};

export function actionTone(action: string): ActionTone {
	const verb = action.split(".").filter(Boolean).at(-1)?.toLowerCase() ?? "";
	return TONE_BY_VERB[verb] ?? "muted";
}

export const ACTION_TONE_CLASS: Record<ActionTone, string> = {
	accent: "text-primary",
	success: "text-success",
	warning: "text-warning",
	danger: "text-destructive",
	muted: "text-muted-foreground",
};
