import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useRouter, useSearch } from "@tanstack/react-router";
import toast from "react-hot-toast";

import type { User } from "#/features/auth/api";
import { AUTH_KEYS, authService, meQuery } from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";
import { clearTempToken, getTempToken } from "../utils/tempToken.ts";

export function useLogin2FA() {
	const navigate = useNavigate();
	const router = useRouter();
	const queryClient = useQueryClient();
	const search = useSearch({ strict: false }) as { redirect?: string };

	return useMutation<ApiResponse<User>, Error, string>({
		mutationKey: AUTH_KEYS.login2FA(),
		mutationFn: (totpCode) => {
			const tempToken = getTempToken();
			if (!tempToken) {
				throw new Error("2FA session expired. Please sign in again.");
			}
			return authService.login2FA({
				temp_token: tempToken,
				totp_code: totpCode,
			});
		},
		onSuccess: async () => {
			clearTempToken();
			await queryClient.fetchQuery({ ...meQuery, staleTime: 0 });
			const target = search.redirect;
			if (target?.startsWith("/") && !target.startsWith("//")) {
				router.history.push(target);
			} else {
				navigate({ to: "/dashboard" });
			}
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
