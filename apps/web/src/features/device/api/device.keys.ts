// Query-key factory for the device flow. Same shape as AUTH_KEYS/ONBOARDING_KEYS:
// a stable root plus one narrow key per unit of state, so hooks can target
// exactly what they own when reading, invalidating, or seeding the cache. The
// approval keys are scoped by user code — two open tabs approving different
// codes must not share a cache entry.
export const DEVICE_KEYS = {
  all: ["device"] as const,
  approval: (userCode: string) =>
    [...DEVICE_KEYS.all, "approval", userCode] as const,
  decide: (userCode: string) =>
    [...DEVICE_KEYS.all, "decide", userCode] as const,
  devices: () => [...DEVICE_KEYS.all, "list"] as const,
  revoke: (id: string) => [...DEVICE_KEYS.all, "revoke", id] as const,
};
