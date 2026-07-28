export const SETTINGS_KEYS = {
	all: ["settings"] as const,
	sessions: () => [...SETTINGS_KEYS.all, "sessions"] as const,
	updateProfile: () => [...SETTINGS_KEYS.all, "update-profile"] as const,
	changePassword: () => [...SETTINGS_KEYS.all, "change-password"] as const,
	enable2FA: () => [...SETTINGS_KEYS.all, "enable-2fa"] as const,
	verify2FA: () => [...SETTINGS_KEYS.all, "verify-2fa"] as const,
	disable2FA: () => [...SETTINGS_KEYS.all, "disable-2fa"] as const,
	revokeSession: () => [...SETTINGS_KEYS.all, "revoke-session"] as const,
	revokeOtherSessions: () =>
		[...SETTINGS_KEYS.all, "revoke-other-sessions"] as const,
};
