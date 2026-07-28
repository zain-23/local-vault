import { createFileRoute, redirect } from "@tanstack/react-router";

import { meQuery } from "#/features/auth/api";
import { ApprovalScreen } from "#/features/device/components/index.ts";
import { useDeviceCodeStore } from "#/features/device/stores/useDeviceCodeStore.ts";

// /device/confirmation — the approval step. The code isn't in the URL; it's held
// in a store set by the submit step, so a refresh drops it and the screen shows
// "no code". Approving links the device to *this* account, so the reviewer must
// be signed in: if they aren't, we bounce to login and back here (a client-side
// round-trip, so the in-memory code survives).
export const Route = createFileRoute("/device/confirmation")({
  beforeLoad: async ({ context, location }) => {
    const user = await context.queryClient.ensureQueryData(meQuery);
    if (!user) {
      throw redirect({
        to: "/auth/login",
        search: { redirect: location.href },
      });
    }
  },
  component: ApprovalPage,
});

function ApprovalPage() {
  const userCode = useDeviceCodeStore((s) => s.userCode);
  return <ApprovalScreen userCode={userCode ?? undefined} />;
}
