import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

import { FocusedLayout } from "#/components/shared";
import { meQuery } from "#/features/auth/api";

// Layout route for the CLI device flow (/device/*). Auth required — approving
// a device ties it to *this* account. Preserve the full URL through login so
// ?user_code=… survives the round-trip.
export const Route = createFileRoute("/device")({
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
	component: DeviceRoute,
});

function DeviceRoute() {
	return (
		<FocusedLayout>
			<Outlet />
		</FocusedLayout>
	);
}
