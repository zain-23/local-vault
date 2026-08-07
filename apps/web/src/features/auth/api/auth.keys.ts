export const AUTH_KEYS = {
	all: ["auth"] as const,
	logout: () => [...AUTH_KEYS.all, "logout"] as const,
	me: () => [...AUTH_KEYS.all, "me"] as const,
};
