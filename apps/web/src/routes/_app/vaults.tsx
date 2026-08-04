import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/vaults")({
	staticData: { breadcrumb: "Vaults" },
	component: VaultsLayout,
});

function VaultsLayout() {
	return <Outlet />;
}
