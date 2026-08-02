import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";

import type { Invite, InviteInput } from "#/features/members/api";
import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ApiResponse } from "#/services/api";
import { useWorkspaceStore } from "#/stores";

// POST /members/invite. On success, refresh pending invites (and the roster in
// case the invitee was already somehow listed — cheap and keeps the UI honest).
export function useInviteMember() {
	const queryClient = useQueryClient();
	const workspaceId = useWorkspaceStore((s) => s.active?.id);

	return useMutation<ApiResponse<Invite>, Error, InviteInput>({
		mutationKey: MEMBER_KEYS.invite(workspaceId ?? ""),
		mutationFn: (input) => {
			if (!workspaceId) throw new Error("No active workspace");
			return memberService.invite(workspaceId, input);
		},
		onSuccess: (res) => {
			toast.success(res.message || "Invite sent");
			if (!workspaceId) return;
			// Prefix-match: refreshes list + invites for this workspace.
			queryClient.invalidateQueries({
				queryKey: MEMBER_KEYS.workspace(workspaceId),
			});
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
