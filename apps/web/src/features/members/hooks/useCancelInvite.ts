import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";

import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ApiResponse } from "#/services/api";
import { useWorkspaceStore } from "#/stores";

// DELETE /members/invites/:id — cancel a pending invite.
export function useCancelInvite() {
	const queryClient = useQueryClient();
	const workspaceId = useWorkspaceStore((s) => s.active?.id);

	return useMutation<ApiResponse<null>, Error, string>({
		mutationKey: MEMBER_KEYS.cancelInvite(workspaceId ?? ""),
		mutationFn: (inviteId) => {
			if (!workspaceId) throw new Error("No active workspace");
			return memberService.cancelInvite(workspaceId, inviteId);
		},
		onSuccess: (res) => {
			toast.success(res.message || "Invite cancelled");
			if (!workspaceId) return;
			queryClient.invalidateQueries({
				queryKey: MEMBER_KEYS.invites(workspaceId),
			});
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
