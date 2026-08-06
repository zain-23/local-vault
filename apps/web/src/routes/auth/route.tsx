import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { meQuery } from "#/features/auth/api";
import { AuthLayout } from "#/features/auth/components";

// Layout route for every /auth/* screen: the split-screen shell + brand panel
// render once here, and each child route supplies only its form via <Outlet/>.
export const Route = createFileRoute("/auth")({
	// Same as /_app: session cookies aren't on the web origin, so SSR would
	// always think you're logged out (and could trap a signed-in refresh here).
	ssr: false,
	beforeLoad: async ({ context }) => {
		const user = await context.queryClient.ensureQueryData(meQuery);
		if (user) throw redirect({ to: user.onboarded ? "/" : "/onboarding" });
	},
	component: AuthLayoutRoute,
});

function AuthLayoutRoute() {
	return (
		<AuthLayout>
			<Outlet />
		</AuthLayout>
	);
}
