import { useQuery } from "@tanstack/react-query";

import { vaultsQuery } from "#/features/vaults/api";
import { useWorkspaceStore } from "#/stores";

export function useVaults() {
	const workspaceId = useWorkspaceStore((s) => s.active?.id);
	return useQuery({
		...vaultsQuery(workspaceId ?? ""),
		enabled: !!workspaceId,
	});
}
