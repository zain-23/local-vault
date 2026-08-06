import { useRouterState } from "@tanstack/react-router";

import type { Crumb } from "#/components/layout/types";
import { readBreadcrumb } from "#/components/layout/utils";
import { useVault } from "#/features/vaults/hooks";

export function useBreadcrumbs(): Crumb[] {
	const matches = useRouterState({ select: (s) => s.matches });
	const vaultId = useRouterState({
		select: (s) => {
			const params = s.matches.at(-1)?.params as
				| Record<string, string | undefined>
				| undefined;
			const id = params?.vaultId;
			return typeof id === "string" ? id : undefined;
		},
	});
	const { data: vault } = useVault(vaultId ?? "");

	const crumbs: Crumb[] = [];
	for (const match of matches) {
		const label = readBreadcrumb(match.staticData);
		if (!label) continue;
		crumbs.push({ label, to: match.pathname });
	}

	if (crumbs.length === 0) return crumbs;

	const last = crumbs.at(-1);
	if (last && vaultId) {
		last.label = vault?.name ?? "Vault";
	}

	return crumbs;
}
