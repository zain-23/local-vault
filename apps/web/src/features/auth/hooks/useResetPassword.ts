import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import toast from "react-hot-toast";

import type { ResetPasswordInput } from "#/features/auth/api";
import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useResetPassword() {
	const navigate = useNavigate();

	return useMutation<ApiResponse<null>, Error, ResetPasswordInput>({
		mutationFn: (input) => authService.resetPassword(input),
		mutationKey: AUTH_KEYS.resetPassword(),
		onSuccess: (res) => {
			toast.success(res.message);
			navigate({ to: "/auth/reset-success" });
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
