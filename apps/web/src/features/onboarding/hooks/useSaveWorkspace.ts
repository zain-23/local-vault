import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";

import type { MemberRole } from "#/features/members/types";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import type { Workspace } from "#/features/onboarding/api";
import { ONBOARDING_KEYS, onboardingService } from "#/features/onboarding/api";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore.ts";
import { getWorkspaceIcon } from "../lib/workspaceIcons.tsx";
import type { WorkspaceValues } from "../schemas/index.ts";

type SavedWorkspace = {
	workspace: Workspace;
	role: MemberRole;
};

// Step 1 — persist the workspace name + icon, resilient to a page refresh. Three
// paths, one hook so the component stays declarative:
//   • no existing workspace  -> create (POST /workspaces/)
//   • name or icon changed   -> update (PUT /workspaces/:id)
//   • both unchanged         -> no request, just advance
// Whichever path runs, the resulting workspace becomes active and the wizard
// moves on. `existing` comes from the prefill query, so a refresh after step 1
// updates that workspace instead of silently creating a duplicate.
export function useSaveWorkspace(
	existing: Workspace | undefined,
	options?: { onSuccess?: () => void },
) {
	const queryClient = useQueryClient();
	const setActive = useWorkspaceStore((s) => s.setActive);

	return useMutation<SavedWorkspace, Error, WorkspaceValues>({
		mutationKey: ONBOARDING_KEYS.saveWorkspace(),
		mutationFn: async ({ name, icon }) => {
			if (!existing) {
				const res = await onboardingService.createWorkspace({ name, icon });
				return {
					workspace: res.data.workspace,
					role: parseMemberRole(res.data.role) ?? "owner",
				};
			}
			const role =
				useWorkspaceStore.getState().active?.role ??
				("owner" satisfies MemberRole);
			// Untouched name and icon: skip the round-trip and reuse what we already have.
			if (
				existing.name === name &&
				getWorkspaceIcon(existing.icon).id === icon
			) {
				return { workspace: existing, role };
			}
			const res = await onboardingService.updateWorkspace(existing.id, {
				name,
				icon,
			});
			return { workspace: res.data, role };
		},
		onSuccess: ({ workspace, role }) => {
			setActive({
				id: workspace.id,
				name: workspace.name,
				icon: workspace.icon,
				role,
			});
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
