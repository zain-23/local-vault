import { queryOptions } from "@tanstack/react-query";

import { ONBOARDING_KEYS, onboardingService } from "#/features/onboarding/api";

// The caller's workspaces. Used to prefill Step 1 so a refreshed onboarding page
// resumes from the workspace already created instead of making a second one.
export const workspacesQuery = queryOptions({
  queryKey: ONBOARDING_KEYS.workspaces(),
  queryFn: async () => {
    const res = await onboardingService.listWorkspaces();
    return res.data;
  },
  staleTime: 5 * 60 * 1000,
});

// The authorized-device list, used by Step 3 to detect when a terminal links.
// Polling cadence lives in the hook (refetchInterval), not here — this stays a
// plain description of "how to fetch the list" so it can be reused elsewhere.
// staleTime 0: while waiting to connect, we always want the freshest answer.
export const devicesQuery = queryOptions({
  queryKey: ONBOARDING_KEYS.devices(),
  queryFn: async () => {
    const res = await onboardingService.listDevices();
    return res.data;
  },
  staleTime: 0,
});
