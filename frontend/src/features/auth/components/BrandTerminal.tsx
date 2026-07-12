import type * as React from "react";

import { cn } from "#/lib/utils.ts";
import type { TerminalEntry } from "../types";

// The brand "signature" session: the product doing its one job — encrypt
// on-device, sync peer-to-peer, inject at runtime. Add/reorder commands here.
const SESSION: TerminalEntry[] = [
	{ id: "init", cmd: "lv init" },
	{
		id: "init-out",
		out: (
			<>
				✓ Created vault <span className="text-foreground">prod-api</span>
				<br />✓ Generated keypair · fingerprint{" "}
				<span className="text-foreground">b427d8d8…</span>
			</>
		),
	},
	{ id: "add", cmd: "lv add STRIPE_SECRET_KEY=sk_live_•••" },
	{ id: "add-out", out: <>✓ Encrypted with AES-256-GCM</> },
	{ id: "push", cmd: "lv push" },
	{
		id: "push-out",
		out: (
			<>
				✓ Synced to 5 peers · <span className="text-success">0 bytes</span> to
				our servers
			</>
		),
	},
	{ id: "inject", cmd: "lv inject --npm start" },
];

// One "traffic light" dot in the terminal title bar.
function Dot({ color }: { color: string }) {
	return (
		<span
			className="size-[7px] rounded-full"
			style={{ background: color }}
			aria-hidden="true"
		/>
	);
}

// A prompt line: `$ <cmd>`.
function Line({ children }: { children: React.ReactNode }) {
	return (
		<div>
			<span className="mr-2 select-none text-primary">$</span>
			<span>{children}</span>
		</div>
	);
}

// Indented command output.
function Output({ children }: { children: React.ReactNode }) {
	return <div className="mb-1.5 pl-4 text-muted-foreground">{children}</div>;
}

function BrandTerminal({ className }: { className?: string }) {
	return (
		<div
			className={cn(
				"overflow-hidden rounded-xl border border-border bg-background shadow-2xl",
				className,
			)}
		>
			<div className="flex h-[34px] items-center gap-[7px] border-b border-border bg-secondary px-3.5">
				<Dot color="#ff5f57" />
				<Dot color="#febc2e" />
				<Dot color="#28c840" />
				<span className="flex-1 text-center font-mono text-[11px] text-muted-foreground">
					lv — secure sync
				</span>
			</div>
			<div className="px-[18px] py-4 font-mono text-[12.5px] leading-[1.85]">
				{SESSION.map((entry) =>
					"cmd" in entry ? (
						<Line key={entry.id}>{entry.cmd}</Line>
					) : (
						<Output key={entry.id}>{entry.out}</Output>
					),
				)}
			</div>
		</div>
	);
}

export { BrandTerminal };
