import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { LoginForm } from "#/features/auth/components/index.ts";
import { seo } from "#/utils/seo.ts";

export const Route = createFileRoute("/auth/login")({
	head: () => seo(PAGE_META["/auth/login"]),
	// `redirect` is where to send the user after a successful login — set by
	// guards (e.g. the device approval screen) that bounced them here. `error`
	// carries an OAuth failure surfaced by the callback.
	validateSearch: (
		search: Record<string, unknown>,
	): { error?: string; redirect?: string } => {
		return {
			...(typeof search.error === "string" ? { error: search.error } : {}),
			...(typeof search.redirect === "string"
				? { redirect: search.redirect }
				: {}),
		};
	},
	component: LoginForm,
});
