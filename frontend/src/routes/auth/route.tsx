import { createFileRoute, Outlet } from "@tanstack/react-router";

import { AuthLayout } from "#/features/auth/components/index.ts";

// Layout route for every /auth/* screen: the split-screen shell + brand panel
// render once here, and each child route supplies only its form via <Outlet/>.
export const Route = createFileRoute("/auth")({
  component: AuthLayoutRoute,
});

function AuthLayoutRoute() {
  return (
    <AuthLayout>
      <Outlet />
    </AuthLayout>
  );
}
