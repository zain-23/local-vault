import { createFileRoute, Outlet } from "@tanstack/react-router";

import { FocusedLayout } from "#/components/shared";

// Layout route for the CLI device flow (/device/*). UI only for now — no auth
// guard; the shared focused card renders once here around each step.
export const Route = createFileRoute("/device")({
  component: DeviceRoute,
});

function DeviceRoute() {
  return (
    <FocusedLayout>
      <Outlet />
    </FocusedLayout>
  );
}
