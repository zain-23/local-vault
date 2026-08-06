import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useSendMagicLink() {
	return useMutation<ApiResponse<null>, Error, string>({
		mutationFn: (email) => authService.sendMagicLink(email),
		mutationKey: AUTH_KEYS.magicLink(),
		onSuccess: (res) => {
			toast.success(res.message);
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
