import type { ListMembersParams } from "./member.types.ts";

// Query-key factory for the members feature. Every key is scoped by workspaceId
// so switching workspaces never serves another team's roster from cache.
export const MEMBER_KEYS = {
  all: ["members"] as const,
  workspace: (workspaceId: string) =>
    [...MEMBER_KEYS.all, workspaceId] as const,
  list: (workspaceId: string, params: ListMembersParams = {}) =>
    [...MEMBER_KEYS.workspace(workspaceId), "list", params] as const,
  invites: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "invites"] as const,
  invite: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "invite"] as const,
  cancelInvite: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "cancel-invite"] as const,
  changeRole: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "change-role"] as const,
  remove: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "remove"] as const,
  join: (workspaceId: string) =>
    [...MEMBER_KEYS.workspace(workspaceId), "join"] as const,
};
