import { createFileRoute } from "@tanstack/react-router";

import { ApprovalScreen } from "#/features/device/components/index.ts";

// /device/confirmation — the approval step where the user reviews the device
// details and approves or denies. The code arrives as ?user_code=WDJF-X4K2 from
// the submit step; the screen is UI-only for now, so it isn't read yet.
export const Route = createFileRoute("/device/confirmation")({
  component: ApprovalPage,
});

function ApprovalPage() {
  return <ApprovalScreen />;
}
