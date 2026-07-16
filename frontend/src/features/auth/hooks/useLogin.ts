import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
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
  const queryClient = useQueryClient();

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
      navigate({ to: "/" });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
