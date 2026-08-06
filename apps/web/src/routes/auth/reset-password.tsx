import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { ResetPasswordForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/reset-password")({
	head: () => seo(PAGE_META["/auth/reset-password"]),
	// Pull ?token=… off the email link.
	validateSearch: (search: Record<string, unknown>): { token?: string } => ({
		token: typeof search.token === "string" ? search.token : undefined,
	}),
	component: RouteComponent,
});

function RouteComponent() {
	const { token } = Route.useSearch();
	return <ResetPasswordForm token={token} />;
}
