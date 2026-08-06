import { Vault } from "lucide-react";
import { motion } from "motion/react";

import { Badge } from "#/components/ui/Badge.tsx";
import { buttonVariants } from "#/components/ui/Button.tsx";
import type { ConsoleStat, VaultRow } from "#/features/marketing/types";
import {
	CONSOLE_NAV_ITEMS,
	CONSOLE_STATS,
	CONSOLE_VAULTS,
} from "#/features/marketing/utils/constants.ts";
import { tableGroup, tableRow } from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";
import { AnimatedCount } from "./AnimatedCount.tsx";

function ConsoleSidebar() {
	return (
		<aside className="hidden border-r border-border bg-sidebar px-2.5 py-3.5 md:block">
			<div className="px-2 pt-1.5 pb-2 text-[11px] font-medium tracking-[0.06em] text-muted-foreground uppercase">
				Workspace
			</div>
			{CONSOLE_NAV_ITEMS.map(({ id, label, icon: Icon, active }) => (
				<div
					key={id}
					className={cn(
						"flex items-center gap-2.5 rounded-md px-2.5 py-1.5 text-[13.5px] text-muted-foreground",
						active && "bg-accent font-medium text-foreground",
					)}
				>
					<Icon className="size-[15px]" />
					{label}
				</div>
			))}
		</aside>
	);
}

function StatCard({ stat }: { stat: ConsoleStat }) {
	return (
		<div className="rounded-lg border border-border bg-background px-3.5 py-3">
			<div className="text-[11.5px] text-muted-foreground">{stat.label}</div>
			<div
				className={cn(
					"mt-0.5 text-[21px] font-semibold tracking-tight tabular-nums",
					stat.tone === "success" && "text-success",
				)}
			>
				{stat.countUp ? <AnimatedCount value={stat.value} /> : stat.value}
			</div>
		</div>
	);
}

function VaultTableRow({ vault }: { vault: VaultRow }) {
	return (
		<motion.tr
			variants={tableRow}
			className="border-b border-border transition-colors last:border-b-0 hover:bg-accent/45"
		>
			<td className="p-2.5">
				<span className="flex items-center gap-2 font-medium">
					<Vault className="size-[15px] text-primary" />
					{vault.name}
				</span>
			</td>
			<td className="hidden p-2.5 font-mono text-xs text-muted-foreground sm:table-cell">
				{vault.secrets}
			</td>
			<td className="hidden p-2.5 font-mono text-xs text-muted-foreground sm:table-cell">
				{vault.members}
			</td>
			<td className="p-2.5 font-mono text-xs text-muted-foreground">
				{vault.lastSync}
			</td>
			<td className="p-2.5">
				<Badge
					className={cn(
						"rounded-md",
						vault.status === "synced"
							? "border-success/35 bg-success/10 text-success"
							: "border-accent-border bg-accent-soft text-primary",
					)}
				>
					{vault.statusLabel}
				</Badge>
			</td>
		</motion.tr>
	);
}

/** The web console half of the product tour: sidebar, live counters, vault list. */
function ConsolePane() {
	return (
		<div className="grid min-h-[330px] grid-cols-1 md:grid-cols-[196px_1fr]">
			<ConsoleSidebar />

			<div className="px-4.5 py-4">
				<div className="mb-3.5 flex items-center justify-between gap-3">
					<h4 className="text-[15px] font-semibold">Vaults</h4>
					{/* A picture of a button — not focusable, it does nothing here. */}
					<span className={buttonVariants({ variant: "outline", size: "sm" })}>
						New vault
					</span>
				</div>

				<div className="mb-4 grid grid-cols-1 gap-2.5 sm:grid-cols-3">
					{CONSOLE_STATS.map((stat) => (
						<StatCard key={stat.id} stat={stat} />
					))}
				</div>

				<table className="w-full border-collapse text-[13px]">
					<thead>
						<tr>
							<th className="border-b border-border p-2.5 text-left text-[11.5px] font-medium whitespace-nowrap text-muted-foreground">
								Vault
							</th>
							<th className="hidden border-b border-border p-2.5 text-left text-[11.5px] font-medium whitespace-nowrap text-muted-foreground sm:table-cell">
								Secrets
							</th>
							<th className="hidden border-b border-border p-2.5 text-left text-[11.5px] font-medium whitespace-nowrap text-muted-foreground sm:table-cell">
								Members
							</th>
							<th className="border-b border-border p-2.5 text-left text-[11.5px] font-medium whitespace-nowrap text-muted-foreground">
								Last sync
							</th>
							<th className="border-b border-border p-2.5 text-left text-[11.5px] font-medium whitespace-nowrap text-muted-foreground">
								Status
							</th>
						</tr>
					</thead>
					<motion.tbody
						initial="hidden"
						animate="visible"
						variants={tableGroup}
					>
						{CONSOLE_VAULTS.map((vault) => (
							<VaultTableRow key={vault.id} vault={vault} />
						))}
					</motion.tbody>
				</table>
			</div>
		</div>
	);
}

export { ConsolePane };
