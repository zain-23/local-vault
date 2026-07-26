import { queryOptions } from "@tanstack/react-query";

import { VAULT_KEYS, vaultService } from "#/features/vaults/api";

export function vaultsQuery(workspaceId: string) {
	return queryOptions({
		queryKey: VAULT_KEYS.list(workspaceId),
		queryFn: async () => {
			const res = await vaultService.list(workspaceId);
			return res.data;
		},
		staleTime: 30_000,
	});
}
