export const SETTINGS_KEYS = {
	all: ["settings"] as const,
	sessions: () => [...SETTINGS_KEYS.all, "sessions"] as const,
	updateProfile: () => [...SETTINGS_KEYS.all, "update-profile"] as const,
	revokeSession: () => [...SETTINGS_KEYS.all, "revoke-session"] as const,
	revokeOtherSessions: () =>
		[...SETTINGS_KEYS.all, "revoke-other-sessions"] as const,
};
