import type { QueryClient } from "@tanstack/react-query";
import { redirect } from "@tanstack/react-router";
import { meQuery, type User } from "#/features/auth/api";

type GuardContext = { queryClient: QueryClient };
type GuardDecision = { to: string } | undefined;

/**
 * Shared shape for every auth-gated route: fetch the session (cookie auth
 * lives on the API origin, so this can only run client-side — hence
 * ssr: false), let `decide` redirect based on the result, then run `after`
 * once the guard passes.
 */
export function authGuard(
	decide: (user: User | null) => GuardDecision,
	after?: (context: GuardContext) => Promise<unknown>,
) {
	return {
		ssr: false as const,
		beforeLoad: async ({ context }: { context: GuardContext }) => {
			const user = await context.queryClient.ensureQueryData(meQuery);
			const result = decide(user);
			if (result) throw redirect(result);
			await after?.(context);
		},
	};
}
