import { createFileRoute } from "@tanstack/react-router";
import { MembersPage } from "#/features/members/components/MembersPage.tsx";

export const Route = createFileRoute("/_app/members")({
	staticData: { breadcrumb: "Members" },
	component: MembersPage,
});
