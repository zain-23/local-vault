import { queryOptions } from "@tanstack/react-query";

import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ListMembersParams } from "#/features/members/api";

// Paginated roster for a workspace. Params are part of the key so distinct
// filter/page combinations never collide in cache. staleTime keeps the table
// snappy while typing filters without hammering the API on every remount.
export function membersQuery(
  workspaceId: string,
  params: ListMembersParams = {},
) {
  return queryOptions({
    queryKey: MEMBER_KEYS.list(workspaceId, params),
    queryFn: async () => {
      const res = await memberService.listMembers(workspaceId, params);
      return res.data;
    },
    staleTime: 30_000,
  });
}

// Pending invites for a workspace. Used by the invites management UI (future)
// and invalidated by invite / cancel mutations so the roster stays honest.
export function invitesQuery(workspaceId: string) {
  return queryOptions({
    queryKey: MEMBER_KEYS.invites(workspaceId),
    queryFn: async () => {
      const res = await memberService.listInvites(workspaceId);
      return res.data;
    },
    staleTime: 30_000,
  });
}
