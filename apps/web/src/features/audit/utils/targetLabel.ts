import type { AuditEvent } from "#/features/audit/api";

const TARGET_TYPE_FALLBACK: Record<string, string> = {
	vault: "a vault",
	workspace: "a workspace",
	user: "a member",
};

/**
 * Human-readable label for an audit event's target. Never returns a raw
 * internal id — falls back to a self-referential or type-generic phrase
 * when the backend didn't supply `target_name`.
 */
export function resolveTargetLabel(event: AuditEvent): string | null {
	const name = event.target_name?.trim();
	if (name) return name;

	if (!event.target_id) return null;

	if (event.actor_id && event.target_id === event.actor_id) {
		return "their own account";
	}

	return TARGET_TYPE_FALLBACK[event.target_type ?? ""] ?? "this item";
}
