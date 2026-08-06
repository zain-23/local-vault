import { useQuery } from "@tanstack/react-query";
import type { ListMembersParams } from "#/features/members/api";
import { membersQuery } from "#/features/members/api";
import { useWorkspaceStore } from "#/stores";

// Paginated members for the active workspace. Filters (page/limit/role/search)
// are part of the query key, so each combination caches independently. Disabled
// until a workspace is selected — nothing to list against.
export function useMembers(params: ListMembersParams = {}) {
	const workspaceId = useWorkspaceStore((s) => s.active?.id);
	return useQuery({
		...membersQuery(workspaceId ?? "", params),
		enabled: !!workspaceId,
	});
}
