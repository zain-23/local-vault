import { queryOptions } from "@tanstack/react-query";
import { AUTH_KEYS, authService } from "#/features/auth/api";
import { ApiError } from "#/services/api";

// Single source of truth for session state. Only a 401 resolves to null — that's
// a real answer ("not logged in"). Anything else throws: a 500 or a dead network
// is not the same as a logged-out user and must not redirect to login.
export const meQuery = queryOptions({
  queryKey: AUTH_KEYS.me(),
  queryFn: async () => {
    try {
      const res = await authService.me();
      return res.data;
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) return null;
      throw err;
    }
  },
  retry: false, // a 401 is an answer, not a failure to retry
  staleTime: 5 * 60 * 1000,
});
