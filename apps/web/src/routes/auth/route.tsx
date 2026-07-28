import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

import { AuthLayout } from "#/features/auth/components";
import { meQuery } from "#/features/auth/api";

// Layout route for every /auth/* screen: the split-screen shell + brand panel
// render once here, and each child route supplies only its form via <Outlet/>.
export const Route = createFileRoute("/auth")({
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
