// Query-key factory for the onboarding flow. Same shape as AUTH_KEYS: a stable
// root plus one narrow key per unit of state, so hooks can target exactly what
// they own when reading, invalidating, or seeding the cache.
export const ONBOARDING_KEYS = {
  all: ["onboarding"] as const,
  workspaces: () => [...ONBOARDING_KEYS.all, "workspaces"] as const,
  saveWorkspace: () => [...ONBOARDING_KEYS.all, "save-workspace"] as const,
  devices: () => [...ONBOARDING_KEYS.all, "devices"] as const,
  complete: () => [...ONBOARDING_KEYS.all, "complete"] as const,
};
