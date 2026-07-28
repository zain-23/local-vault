import { useQuery } from "@tanstack/react-query";

import { invitesQuery } from "#/features/members/api";
import { canManageInvites } from "#/features/members/utils/canManageInvites.ts";
import { useWorkspaceStore } from "#/stores";

// Pending invites for the active workspace. Only fetched for owner/admin —
// members must never hit GET /invites (server 403s; we skip the round-trip).
export function useInvites() {
  const workspace = useWorkspaceStore((s) => s.active);
  const workspaceId = workspace?.id;
  const enabled = !!workspaceId && canManageInvites(workspace?.role);

  return useQuery({
    ...invitesQuery(workspaceId ?? ""),
    enabled,
  });
}
