import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useRouter, useSearch } from "@tanstack/react-router";
import { toast } from "sonner";
import type { LoginInput, LoginResult } from "#/features/auth/api";
import {
	AUTH_KEYS,
	authService,
	isTwoFactorRequired,
	meQuery,
} from "#/features/auth/api";
import type { ApiResponse } from "#/services/api";

export function useLogin() {
	const navigate = useNavigate();
	const router = useRouter();
	const queryClient = useQueryClient();
	// Read loosely: this hook is only mounted on /auth/login, but strict:false
	// keeps it decoupled from that route's typed search.
	const search = useSearch({ strict: false }) as { redirect?: string };

	return useMutation<ApiResponse<LoginResult>, Error, LoginInput>({
		mutationFn: (input) => authService.login(input),
		mutationKey: AUTH_KEYS.login(),
		onSuccess: async (res) => {
			if (isTwoFactorRequired(res.data)) {
				navigate({ to: "/auth/two-factor" });
				return;
			}
			// SSR rendered this page logged-out, so the cache holds me:null — and a
			// hydrated query has no queryFn, so invalidating it can't refetch. Pass
			// meQuery's own options to force the call, and await it: the /_app guard
			// reads this cache and would bounce us back to login on a stale null.
			await queryClient.fetchQuery({ ...meQuery, staleTime: 0 });
			// Honor a guard's ?redirect=… back to where the user was headed. Only
			// internal paths ("/…", never "//host") — a rogue link must not turn login
			// into an open redirect. history.push takes the raw path+query as-is.
			const target = search.redirect;
			if (target?.startsWith("/") && !target.startsWith("//")) {
				router.history.push(target);
			} else {
				navigate({ to: "/" });
			}
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
