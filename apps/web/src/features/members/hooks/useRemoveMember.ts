import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { MEMBER_KEYS, memberService } from "#/features/members/api";
import type { ApiResponse } from "#/services/api";
import { useWorkspaceStore } from "#/stores";

// DELETE /members/:userId — remove a member from the active workspace.
export function useRemoveMember() {
  const queryClient = useQueryClient();
  const workspaceId = useWorkspaceStore((s) => s.active?.id);

  return useMutation<ApiResponse<null>, Error, string>({
    mutationKey: MEMBER_KEYS.remove(workspaceId ?? ""),
    mutationFn: (userId) => {
      if (!workspaceId) throw new Error("No active workspace");
      return memberService.removeMember(workspaceId, userId);
    },
    onSuccess: (res) => {
      toast.success(res.message || "Member removed");
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
