import { createFileRoute } from "@tanstack/react-router";

import { VaultDetailPage } from "#/features/vaults/components/VaultDetailPage.tsx";

export const Route = createFileRoute("/_app/vaults/$vaultId")({
  component: VaultDetailRoute,
});

function VaultDetailRoute() {
  const { vaultId } = Route.useParams();
  return <VaultDetailPage vaultId={vaultId} />;
}
