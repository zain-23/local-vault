import { createFileRoute } from "@tanstack/react-router";

import { CheckEmailPanel } from "#/features/auth/components/index.ts";

// The confirmation copy + "use a different address" target depend on which flow
// sent the user here, carried in the URL so the page is refresh-safe.
type CheckEmailSearch = { variant: "reset" | "magic"; email?: string };

export const Route = createFileRoute("/auth/check-email")({
	validateSearch: (search: Record<string, unknown>): CheckEmailSearch => ({
		variant: search.variant === "magic" ? "magic" : "reset",
		email: typeof search.email === "string" ? search.email : undefined,
	}),
	component: CheckEmail,
});

function CheckEmail() {
	const { variant, email } = Route.useSearch();
	return <CheckEmailPanel variant={variant} email={email} />;
}
