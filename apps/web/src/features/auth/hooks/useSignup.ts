import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

import type { SignupInput } from "#/features/auth/api";
import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useSignup() {
	return useMutation<ApiResponse<string>, Error, SignupInput>({
		mutationFn: (input) => authService.signup(input),
		mutationKey: AUTH_KEYS.signup(),
		// Failures throw from the api client — this only runs on a real signup.
		onSuccess: (res) => {
			toast.success(res.message);
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
