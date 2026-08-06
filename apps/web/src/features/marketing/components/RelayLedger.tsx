import { RELAY_LEDGER } from "#/features/marketing/utils/constants.ts";
import { cn } from "#/lib/utils.ts";

// Full disclosure of what the relay holds while forwarding a blob. The dimmed
// rows are the point: they're the fields it never receives at all.
function RelayLedger() {
	return (
		<div className="mt-4 overflow-hidden rounded-lg border border-border">
			<div className="flex justify-between border-b border-border bg-muted px-3 py-2 text-[11px] font-medium tracking-[0.06em] text-muted-foreground uppercase">
				<span>What the relay can see</span>
				<span>Value</span>
			</div>
			{RELAY_LEDGER.map((entry) => (
				<div
					key={entry.id}
					className="flex justify-between gap-3.5 border-b border-border px-3 py-2 font-mono text-xs last:border-b-0"
				>
					<span className="text-muted-foreground">{entry.key}</span>
					<span
						className={cn(
							entry.exposed
								? "text-primary"
								: "text-muted-foreground opacity-60",
						)}
					>
						{entry.value}
					</span>
				</div>
			))}
		</div>
	);
}

export { RelayLedger };
