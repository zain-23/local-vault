import { createFileRoute } from "@tanstack/react-router";

import { VerifyEmailPanel } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/verify-email")({
	// Pull ?token=… off the email link; undefined when the user just landed here post-signup.
	validateSearch: (search: Record<string, unknown>): { token?: string } => ({
		token: typeof search.token === "string" ? search.token : undefined,
	}),
	component: RouteComponent,
});

function RouteComponent() {
	const { token } = Route.useSearch();
	return <VerifyEmailPanel token={token} />;
}
