import { useMutation } from "@tanstack/react-query";

import { AUTH_KEYS, authService } from "#/features/auth/api";

// Verifies the email token. No toast/redirect here — the panel renders the
// pending / success / error states inline. The client returns { success:false }
// on a bad token (it doesn't throw), so we throw here to flip the mutation to
// error state and surface the server message via `error`.
export function useVerifyEmail() {
  return useMutation<string, Error, string>({
    mutationKey: AUTH_KEYS.verifyEmail(),
    mutationFn: async (token) => {
      const res = await authService.verifyEmail(token);
      if (!res.success) {
        throw new Error(res.message);
      }
      return res.message; // e.g. "email verified"
    },
  });
}
