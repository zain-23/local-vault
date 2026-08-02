import { createFileRoute } from "@tanstack/react-router";

import { SubmitCodeForm } from "#/features/device/components/index.ts";
import { normalizeUserCode } from "#/features/device/schemas/index.ts";

// /device — the submit-code step. The terminal opens this URL and may append the
// code (?user_code=WDJF-X4K2) to prefill it; we normalise whatever arrives.
export const Route = createFileRoute("/device/")({
	validateSearch: (search: Record<string, unknown>): { user_code?: string } => {
		const raw = search.user_code;
		return typeof raw === "string" ? { user_code: raw } : {};
	},
	component: SubmitCodePage,
});

function SubmitCodePage() {
	const { user_code } = Route.useSearch();
	return <SubmitCodeForm defaultCode={normalizeUserCode(user_code ?? "")} />;
}
