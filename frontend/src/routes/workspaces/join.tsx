import { createFileRoute, redirect } from "@tanstack/react-router";

import { meQuery } from "#/features/auth/api";
import { JoinWorkspaceScreen } from "#/features/members/components/JoinWorkspaceScreen.tsx";

// /workspaces/join?token=&workspaceId= — email invite landing page.
// Auth required (join reads the caller from the JWT). Preserve the full URL
// through login so the token survives the round-trip.
export const Route = createFileRoute("/workspaces/join")({
  validateSearch: (
    search: Record<string, unknown>,
  ): { token?: string; workspaceId?: string } => ({
    token: typeof search.token === "string" ? search.token : undefined,
    workspaceId:
      typeof search.workspaceId === "string" ? search.workspaceId : undefined,
  }),
  beforeLoad: async ({ context, location }) => {
    const user = await context.queryClient.ensureQueryData(meQuery);
    if (!user) {
      throw redirect({
        to: "/auth/login",
        search: { redirect: location.href },
      });
    }
  },
  component: JoinPage,
});

function JoinPage() {
  const { token, workspaceId } = Route.useSearch();
  return <JoinWorkspaceScreen token={token} workspaceId={workspaceId} />;
}
