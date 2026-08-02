import type { AssignableRole, MemberRole } from "#/features/members/types";
import { roleLabel } from "./roleLabel.ts";

// All workspace roles. Source of truth for filters, badges, and URL parsers.
export const MEMBER_ROLES = [
	"owner",
	"admin",
	"member",
] as const satisfies readonly MemberRole[];

// The role values a member can be invited as / changed to. The owner role is
// never assignable through the UI (server rejects it too).
export const ASSIGNABLE_ROLES = [
	"admin",
	"member",
] as const satisfies readonly AssignableRole[];

// Role filter sentinel for "no role filter" in the URL / Select.
export const ROLE_FILTER_ALL = "all" as const;

export type RoleFilterValue = MemberRole | typeof ROLE_FILTER_ALL;

// Display labels for the role filter Select (All + every MemberRole).
export const ROLE_FILTER_LABELS = new Map<RoleFilterValue, string>([
	[ROLE_FILTER_ALL, "All roles"],
	...MEMBER_ROLES.map((role) => [role, roleLabel(role)] as const),
]);
