import { createFileRoute } from "@tanstack/react-router";

import { JoinWorkspaceScreen } from "#/features/members/components/JoinWorkspaceScreen.tsx";

// /workspaces/join?token=&workspaceId= — email invite landing page.
// Auth is on the /workspaces layout.
export const Route = createFileRoute("/workspaces/join")({
	validateSearch: (
		search: Record<string, unknown>,
	): { token?: string; workspaceId?: string } => ({
		token: typeof search.token === "string" ? search.token : undefined,
		workspaceId:
			typeof search.workspaceId === "string" ? search.workspaceId : undefined,
	}),
	component: JoinPage,
});

function JoinPage() {
	const { token, workspaceId } = Route.useSearch();
	return <JoinWorkspaceScreen token={token} workspaceId={workspaceId} />;
}
