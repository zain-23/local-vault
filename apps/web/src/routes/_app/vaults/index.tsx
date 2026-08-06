import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { VaultsPage } from "#/features/vaults/components/VaultsPage.tsx";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/vaults/")({
	head: () => seo(PAGE_META["/_app/vaults/"]),
	component: VaultsPage,
});
