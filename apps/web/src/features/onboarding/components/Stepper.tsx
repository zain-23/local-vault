import { cn } from "#/lib/utils.ts";

// Reusable, presentational progress indicator for any multi-step flow.
// Knows nothing about onboarding — just draws "STEP x OF n" + segment bars.
function Stepper({ current, total = 4 }: { current: number; total?: number }) {
	const pct = Math.round((current / total) * 100);

	return (
		<div className="mb-7">
			<div className="mb-2 flex justify-between font-mono text-xs text-muted-foreground">
				<span>
					STEP {current} OF {total}
				</span>
				<span>{pct}%</span>
			</div>

			<div className="flex gap-1">
				{/* one bar per step (1..total); filled with the accent once passed */}
				{Array.from({ length: total }, (_, i) => i + 1).map((n) => (
					<div
						key={n}
						className={cn(
							"h-[3px] flex-1 rounded-full transition-colors",
							n <= current ? "bg-primary" : "bg-border",
						)}
					/>
				))}
			</div>
		</div>
	);
}

export { Stepper };
