import { useMutation } from "@tanstack/react-query";

import { AUTH_KEYS, authService } from "#/features/auth/api";

// Verifies the email token. No toast/redirect here — the panel renders the
// pending / success / error states inline. A bad token throws from the api
// client, which flips the mutation to error state and surfaces the server
// message via `error`.
export function useVerifyEmail() {
	return useMutation<string, Error, string>({
		mutationKey: AUTH_KEYS.verifyEmail(),
		mutationFn: async (token) => {
			const res = await authService.verifyEmail(token);
			return res.message; // e.g. "email verified"
		},
	});
}
