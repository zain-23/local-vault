import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import toast from "react-hot-toast";

import { meQuery } from "#/features/auth/api";
import type { JoinInput, JoinResult } from "#/features/members/api";
import { MEMBER_KEYS, memberService } from "#/features/members/api";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import { ONBOARDING_KEYS, workspacesQuery } from "#/features/onboarding/api";
import type { ApiResponse } from "#/services/api";
import { useWorkspaceStore } from "#/stores";

export interface JoinWorkspaceVariables extends JoinInput {
	workspaceId: string;
}

// POST /members/join — accept an invite token. The joiner isn't a member yet, so
// workspaceId comes from the caller (invite link / join page), not the store.
// On success: refresh workspaces, set the joined workspace active, then route
// into the app (or onboarding if the account isn't finished yet).
export function useJoinWorkspace() {
	const queryClient = useQueryClient();
	const navigate = useNavigate();
	const setActive = useWorkspaceStore((s) => s.setActive);

	return useMutation<ApiResponse<JoinResult>, Error, JoinWorkspaceVariables>({
		mutationKey: MEMBER_KEYS.join(""),
		mutationFn: ({ workspaceId, token }) =>
			memberService.join(workspaceId, { token }),
		onSuccess: async (res, { workspaceId }) => {
			toast.success(res.message || "Joined workspace");

			queryClient.invalidateQueries({
				queryKey: MEMBER_KEYS.workspace(workspaceId),
			});

			const list = await queryClient.fetchQuery(workspacesQuery);
			const match = list?.find((w) => w.workspace.id === workspaceId);
			const role = parseMemberRole(match?.role ?? res.data.role) ?? "member";

			setActive({
				id: workspaceId,
				name: match?.workspace.name ?? "Workspace",
				plan: "Free",
				role,
			});

			const user = queryClient.getQueryData(meQuery.queryKey);
			await navigate({
				to: user?.onboarded ? "/members" : "/onboarding",
			});
		},
		meta: { invalidates: [ONBOARDING_KEYS.workspaces()] },
	});
}
