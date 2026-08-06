import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { MembersPage } from "#/features/members/components/MembersPage.tsx";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/_app/members")({
	head: () => seo(PAGE_META["/_app/members"]),
	staticData: { breadcrumb: "Members" },
	component: MembersPage,
});
