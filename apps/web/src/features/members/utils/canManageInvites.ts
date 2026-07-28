import type { MemberRole } from "#/features/members/types";
import { MEMBER_ROLES } from "./constants.ts";

// Owner/admin capabilities shared by invites + member row actions (remove).
// Matches RequireRole(owner, admin) on invite routes and DELETE /members/:userId.
export function canManageInvites(role: MemberRole | null | undefined): boolean {
  return role === "owner" || role === "admin";
}

// Change-role is owner-only — PUT /members/:userId/role.
export function canChangeMemberRole(
  role: MemberRole | null | undefined,
): boolean {
  return role === "owner";
}

export function parseMemberRole(
  role: string | null | undefined,
): MemberRole | null {
  if (role && (MEMBER_ROLES as readonly string[]).includes(role)) {
    return role as MemberRole;
  }
  return null;
}
