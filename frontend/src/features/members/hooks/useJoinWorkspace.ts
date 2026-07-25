import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import type { JoinInput, JoinResult } from "#/features/members/api";
import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ApiResponse } from "#/services/api";

export interface JoinWorkspaceVariables extends JoinInput {
  workspaceId: string;
}

// POST /members/join — accept an invite token. The joiner isn't a member yet, so
// workspaceId comes from the caller (invite link / join page), not the store.
export function useJoinWorkspace() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse<JoinResult>, Error, JoinWorkspaceVariables>({
    mutationKey: MEMBER_KEYS.join(""),
    mutationFn: ({ workspaceId, token }) =>
      memberService.join(workspaceId, { token }),
    onSuccess: (res, { workspaceId }) => {
      toast.success(res.message || "Joined workspace");
      queryClient.invalidateQueries({
        queryKey: MEMBER_KEYS.workspace(workspaceId),
      });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
