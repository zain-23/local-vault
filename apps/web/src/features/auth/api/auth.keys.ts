export const AUTH_KEYS = {
	all: ["auth"] as const,
	login: () => [...AUTH_KEYS.all, "login"] as const,
	signup: () => [...AUTH_KEYS.all, "signup"] as const,
	forgotPassword: () => [...AUTH_KEYS.all, "forgot-password"] as const,
	resetPassword: () => [...AUTH_KEYS.all, "reset-password"] as const,
	magicLink: () => [...AUTH_KEYS.all, "magic-link"] as const,
	verifyEmail: () => [...AUTH_KEYS.all, "verify-email"] as const,
	me: () => [...AUTH_KEYS.all, "me"] as const,
};
