import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AuthLayout } from "#/features/auth/components";
import { authGuard } from "#/features/auth/utils/authGuard";

// Layout route for every /auth/* screen: the split-screen shell + brand panel
// render once here, and each child route supplies only its form via <Outlet/>.
export const Route = createFileRoute("/auth")({
	...authGuard((user) => {
		if (user) return { to: user.onboarded ? "/dashboard" : "/onboarding" };
	}),
	component: AuthLayoutRoute,
});

function AuthLayoutRoute() {
	return (
		<AuthLayout>
			<Outlet />
		</AuthLayout>
	);
}
