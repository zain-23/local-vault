import { createFileRoute } from "@tanstack/react-router";
import { VaultsPage } from "#/features/vaults/components/VaultsPage.tsx";

export const Route = createFileRoute("/_app/vaults/")({
	component: VaultsPage,
});
