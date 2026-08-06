import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";

import type { AssignableRole, Member } from "#/features/members/api";
import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ApiResponse } from "#/services/api";
import { useWorkspaceStore } from "#/stores";

export interface ChangeRoleVariables {
	userId: string;
	role: AssignableRole;
}

// PUT /members/:userId/role. Invalidates every list cache for this workspace so
// filtered pages refresh with the new role.
export function useChangeRole() {
	const queryClient = useQueryClient();
	const workspaceId = useWorkspaceStore((s) => s.active?.id);

	return useMutation<ApiResponse<Member>, Error, ChangeRoleVariables>({
		mutationKey: MEMBER_KEYS.changeRole(workspaceId ?? ""),
		mutationFn: ({ userId, role }) => {
			if (!workspaceId) throw new Error("No active workspace");
			return memberService.changeRole(workspaceId, userId, { role });
		},
		onSuccess: (res) => {
			toast.success(res.message || "Role updated");
			if (!workspaceId) return;
			queryClient.invalidateQueries({
				queryKey: MEMBER_KEYS.workspace(workspaceId),
			});
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
