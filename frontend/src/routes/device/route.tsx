import { createFileRoute, Outlet } from "@tanstack/react-router";

import { DeviceLayout } from "#/features/device/components/index.ts";

// Layout route for the CLI device flow (/device/*). UI only for now — no auth
// guard; the shared focused card renders once here around each step.
export const Route = createFileRoute("/device")({
  component: DeviceRoute,
});

function DeviceRoute() {
  return (
    <DeviceLayout>
      <Outlet />
    </DeviceLayout>
  );
}
