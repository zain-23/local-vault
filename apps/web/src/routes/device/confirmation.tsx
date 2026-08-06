import { createFileRoute } from "@tanstack/react-router";
import { PAGE_META } from "#/constants";
import { ApprovalScreen } from "#/features/device/components/index.ts";
import { useDeviceCodeStore } from "#/features/device/stores/useDeviceCodeStore.ts";
import { seo } from "#/utils/seo.ts";

// /device/confirmation — the approval step. Auth is on the /device layout.
// The code isn't in the URL; it's held in a store set by the submit step, so a
// refresh drops it and the screen shows "no code".
export const Route = createFileRoute("/device/confirmation")({
	head: () => seo(PAGE_META["/device/confirmation"]),
	component: ApprovalPage,
});

function ApprovalPage() {
	const userCode = useDeviceCodeStore((s) => s.userCode);
	return <ApprovalScreen userCode={userCode ?? undefined} />;
}
