import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { VaultDetailPage } from "#/features/vaults/components/VaultDetailPage.tsx";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/vaults/$vaultId")({
	head: () => seo(PAGE_META["/_app/vaults/$vaultId"]),
	staticData: { breadcrumb: "Vault" },
	component: VaultDetailRoute,
});

function VaultDetailRoute() {
	const { vaultId } = Route.useParams();
	return <VaultDetailPage vaultId={vaultId} />;
}
