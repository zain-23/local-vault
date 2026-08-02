import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

import { FocusedLayout } from "#/components/shared";
import { meQuery } from "#/features/auth/api";

// Focused card shell for workspace invite flows (/workspaces/*). Auth required
// (join reads the caller from the JWT). Preserve the full URL through login so
// the invite token survives the round-trip.
export const Route = createFileRoute("/workspaces")({
	ssr: false,
	beforeLoad: async ({ context, location }) => {
		const user = await context.queryClient.ensureQueryData(meQuery);
		if (!user) {
			throw redirect({
				to: "/auth/login",
				search: { redirect: location.href },
			});
		}
	},
	component: WorkspacesRoute,
});

function WorkspacesRoute() {
	return (
		<FocusedLayout>
			<Outlet />
		</FocusedLayout>
	);
}
