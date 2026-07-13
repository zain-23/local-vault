import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";

import type { ApiResponse } from "#/services/api";
import { AUTH_KEYS, authService } from "#/features/auth/api";
import type { LoginInput, LoginResult } from "#/features/auth/api";

export function useLogin() {
  const navigate = useNavigate();

  return useMutation<ApiResponse<LoginResult>, Error, LoginInput>({
    mutationFn: (input) => authService.login(input),
    mutationKey: AUTH_KEYS.login(),
    onSuccess: (res) => {
      if (!res.success) {
        throw new Error(res.message);
      }

      if ("requires_2fa" in res.data) {
        navigate({ to: "/auth/two-factor" });
        return;
      }
      navigate({ to: "/" });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
}
