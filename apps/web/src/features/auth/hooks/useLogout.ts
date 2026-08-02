import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import toast from "react-hot-toast";

import { AUTH_KEYS, authService, meQuery } from "#/features/auth/api";

export function useLogout() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	return useMutation({
		mutationKey: AUTH_KEYS.logout(),
		mutationFn: () => authService.logout(),
		onSuccess: async () => {
			// Drop the cached user so auth guards don't treat us as still signed in.
			queryClient.setQueryData(meQuery.queryKey, null);
			await queryClient.resetQueries({ queryKey: AUTH_KEYS.all });
			navigate({ to: "/auth/login" });
		},
		onError: (error: Error) => {
			toast.error(error.message);
		},
	});
}
