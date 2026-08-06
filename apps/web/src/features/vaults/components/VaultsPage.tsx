import { useNavigate } from "@tanstack/react-router";

import { DataTable } from "#/components/ui";
import { useVaults } from "#/features/vaults/hooks";
import { vaultColumns } from "./columns.tsx";

export function VaultsPage() {
	const navigate = useNavigate();
	const { data, isLoading, isError, isFetching } = useVaults();
	const vaults = data ?? [];

	return (
		<div className="flex flex-col gap-6 p-6">
			<div>
				<h1 className="text-xl font-semibold tracking-tight">Vaults</h1>
				<p className="text-sm text-muted-foreground">
					Manage the vaults in your workspace.
				</p>
			</div>

			<DataTable
				columns={vaultColumns}
				data={vaults}
				isLoading={isLoading || isFetching}
				errorMessage={
					isError
						? "Could not load vaults. Check your connection and try again."
						: undefined
				}
				emptyMessage="No vaults yet."
				onRowClick={(row) =>
					navigate({ to: "/vaults/$vaultId", params: { vaultId: row.id } })
				}
			/>
		</div>
	);
}
