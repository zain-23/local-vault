import type { MemberRole } from "#/features/members/types";

export function roleLabel(role: MemberRole): string {
	return role.charAt(0).toUpperCase() + role.slice(1);
}
