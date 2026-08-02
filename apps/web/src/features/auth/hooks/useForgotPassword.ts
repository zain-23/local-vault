import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";

import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useForgotPassword() {
	return useMutation<ApiResponse<null>, Error, string>({
		mutationFn: (email) => authService.forgotPassword(email),
		mutationKey: AUTH_KEYS.forgotPassword(),
		onSuccess: (res) => {
			toast.success(res.message);
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
