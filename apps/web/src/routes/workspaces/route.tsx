import { createFileRoute, Outlet } from "@tanstack/react-router";

import { FocusedLayout } from "#/components/shared";

// Focused card shell for workspace invite flows (/workspaces/*). Same chrome
// as other out-of-app flows — centered card, grid backdrop, no app sidebar.
export const Route = createFileRoute("/workspaces")({
  component: WorkspacesRoute,
});

function WorkspacesRoute() {
  return (
    <FocusedLayout>
      <Outlet />
    </FocusedLayout>
  );
}
