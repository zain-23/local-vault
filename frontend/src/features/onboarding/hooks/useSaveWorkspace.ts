import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import type { Workspace } from "#/features/onboarding/api";
import { ONBOARDING_KEYS, onboardingService } from "#/features/onboarding/api";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore.ts";
import type { WorkspaceValues } from "../schemas/index.ts";

// Step 1 — persist the workspace name, resilient to a page refresh. Three paths,
// one hook so the component stays declarative:
//   • no existing workspace  -> create (POST /workspaces/)
//   • name changed           -> update (PUT /workspaces/:id)
//   • name unchanged         -> no request, just advance
// Whichever path runs, the resulting workspace becomes active and the wizard
// moves on. `existing` comes from the prefill query, so a refresh after step 1
// updates that workspace instead of silently creating a duplicate.
export function useSaveWorkspace(
  existing: Workspace | undefined,
  options?: { onSuccess?: () => void },
) {
  const queryClient = useQueryClient();
  const setActive = useWorkspaceStore((s) => s.setActive);

  return useMutation<Workspace, Error, WorkspaceValues>({
    mutationKey: ONBOARDING_KEYS.saveWorkspace(),
    mutationFn: async ({ name }) => {
      if (!existing) {
        const res = await onboardingService.createWorkspace({ name });
        return res.data.workspace;
      }
      // Untouched name: skip the round-trip and reuse what we already have.
      if (existing.name === name) return existing;
      const res = await onboardingService.updateWorkspace(existing.id, {
        name,
      });
      return res.data;
    },
    onSuccess: (workspace) => {
      // The store carries a `plan` the server doesn't return yet; a workspace
      // starts on Free until billing is wired.
      setActive({ id: workspace.id, name: workspace.name, plan: "Free" });
      // Keep the prefill in sync after a create/rename, and refresh the device
      // list since a linked terminal is scoped to this workspace.
      queryClient.invalidateQueries({ queryKey: ONBOARDING_KEYS.workspaces() });
      queryClient.invalidateQueries({ queryKey: ONBOARDING_KEYS.devices() });
      options?.onSuccess?.();
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
