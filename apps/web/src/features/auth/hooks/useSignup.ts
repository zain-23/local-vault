import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";

import type { SignupInput } from "#/features/auth/api";
import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useSignup() {
	const navigate = useNavigate();

	return useMutation<ApiResponse<string>, Error, SignupInput>({
		mutationFn: (input) => authService.signup(input),
		mutationKey: AUTH_KEYS.signup(),
		// Failures throw from the api client — this only runs on a real signup.
		onSuccess: (res) => {
			toast.success(res.message);
			navigate({ to: "/auth/verify-email" });
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
