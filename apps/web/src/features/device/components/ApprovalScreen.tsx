import { Link } from "@tanstack/react-router";
import { Loader2, Terminal } from "lucide-react";
import { FocusedCard } from "#/components/shared";
import { Button, ErrorMessage, SuccessMessage } from "#/components/ui";
import { useApprovalDetails, useDecideDevice } from "#/features/device/hooks";
import { ApiError } from "#/services/api";

// One key/value line inside the device summary panel.
function DetailRow({
	label,
	children,
}: {
	label: string;
	children: React.ReactNode;
}) {
	return (
		<div className="flex items-center justify-between gap-4 text-[13px]">
			<span className="text-muted-foreground">{label}</span>
			<span className="text-right">{children}</span>
		</div>
	);
}

function Header() {
	return (
		<div className="mb-5 flex items-start gap-3.5">
			<div className="flex size-11 shrink-0 items-center justify-center rounded-xl border border-border bg-muted text-foreground">
				<Terminal className="size-5" />
			</div>
			<div className="min-w-0">
				<h1 className="text-[17px] font-semibold tracking-[-0.01em]">
					Authorize CLI device
				</h1>
				<p className="mt-0.5 text-[13px] text-muted-foreground">
					A new terminal wants to link to your account.
				</p>
			</div>
		</div>
	);
}
// A device links to the whole account here — the workspace it acts in is chosen
// later in the CLI, not on this screen. The route guarantees the user is signed
// in before we get here, so a 401 isn't a case we render.
function ApprovalScreen({ userCode }: { userCode?: string }) {
	const { data, isPending, error } = useApprovalDetails(userCode);
	const decide = useDecideDevice(userCode ?? "");

	// No code in the store — the page was refreshed or opened directly. The code
	// is never in the URL, so there's nothing to recover; send them back to
	// re-enter the one from their terminal.
	if (!userCode) {
		return (
			<FocusedCard header={<Header />}>
				<ErrorMessage title="No device code — enter the code from your terminal again." />
				<Button asChild variant="outline" className="mt-5 w-full">
					<Link to="/device">Enter code</Link>
				</Button>
			</FocusedCard>
		);
	}

	if (isPending) {
		return (
			<FocusedCard header={<Header />}>
				<div className="flex items-center justify-center gap-2 py-6 text-[13px] text-muted-foreground">
					<Loader2 className="size-4 animate-spin" />
					Loading request…
				</div>
			</FocusedCard>
		);
	}

	if (error) {
		// 404 = unknown or expired code; anything else, show the server's words.
		const notFound = error instanceof ApiError && error.status === 404;
		return (
			<FocusedCard header={<Header />}>
				<ErrorMessage
					title={
						notFound ? "This code is invalid or has expired." : error.message
					}
				/>
			</FocusedCard>
		);
	}

	// The decision made this session (undefined until the user acts).
	const decided = decide.isSuccess ? decide.variables : undefined;
	const approved = decided === "approve" || data.status === "approved";
	const denied = decided === "deny" || data.status === "denied";
	const expired = data.status === "expired";
	const settled = approved || denied || expired;

	console.log({ decided, approved, denied, expired, settled });
	return (
		<FocusedCard header={<Header />}>
			<div className="flex flex-col gap-2.5 rounded-xl border border-border bg-muted/40 px-4 py-3.5">
				<DetailRow label="Device name">
					<span className="font-mono">{data.device_name}</span>
				</DetailRow>
				<DetailRow label="IP address">
					<span className="font-mono">{data.ip}</span>
				</DetailRow>
			</div>

			{!settled && (
				<>
					<div className="mt-6 flex gap-3">
						<Button
							variant="outline"
							className="flex-1"
							disabled={decide.isPending}
							onClick={() => decide.mutate("deny")}
						>
							Deny
						</Button>
						<Button
							className="flex-1"
							isLoading={decide.isPending}
							onClick={() => decide.mutate("approve")}
						>
							Approve
						</Button>
					</div>
					<p className="mt-4 text-center text-xs text-muted-foreground">
						Only approve if you started this from your own terminal.
					</p>
				</>
			)}

			{approved && (
				<div className="mt-6">
					<SuccessMessage title="Device linked — you can return to your terminal." />
				</div>
			)}
			{denied && (
				<div className="mt-6">
					<ErrorMessage title="Request denied." />
				</div>
			)}
			{expired && !approved && !denied && (
				<div className="mt-6">
					<ErrorMessage title="This request has expired." />
				</div>
			)}
		</FocusedCard>
	);
}

export { ApprovalScreen };
