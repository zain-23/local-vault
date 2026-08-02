import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

import { AUTH_KEYS, authService } from "#/features/auth/api";

// Verifies the email token. Success/error surface as toasts; the panel
// still uses mutation state for the continue/back buttons.
export function useVerifyEmail() {
	return useMutation<string, Error, string>({
		mutationKey: AUTH_KEYS.verifyEmail(),
		mutationFn: async (token) => {
			const res = await authService.verifyEmail(token);
			return res.message; // e.g. "email verified"
		},
		onSuccess: (message) => {
			toast.success(message);
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
