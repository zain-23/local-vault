import { Check, Copy } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";

import { useCopyToClipboard } from "#/features/marketing/hooks";
import { INSTALL_COMMAND } from "#/features/marketing/utils/constants.ts";
import { SNAP_TRANSITION } from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";

// `$ brew install localvault` with a copy button that confirms itself.
function InstallCommand({ className }: { className?: string }) {
	const { copied, copy } = useCopyToClipboard(INSTALL_COMMAND);

	return (
		<div
			className={cn(
				"inline-flex max-w-full items-center gap-3 rounded-lg border border-border bg-card py-2 pr-2 pl-3.5 font-mono text-[13.5px]",
				className,
			)}
		>
			<span className="text-muted-foreground select-none">$</span>
			<span className="truncate text-foreground">{INSTALL_COMMAND}</span>
			<button
				type="button"
				onClick={copy}
				aria-label={`Copy "${INSTALL_COMMAND}" to the clipboard`}
				className={cn(
					"inline-flex h-7 flex-none items-center justify-center gap-1.5 rounded-md border border-border bg-secondary px-2.5 font-sans text-xs font-medium text-muted-foreground transition-colors",
					"hover:border-border-strong hover:text-foreground",
					copied && "border-success/40 text-success hover:text-success",
				)}
			>
				{/* Icons swap in place — same box, so the button never reflows. */}
				<span className="relative inline-flex size-3 items-center justify-center">
					<AnimatePresence initial={false} mode="popLayout">
						{copied ? (
							<motion.span
								key="done"
								initial={{ opacity: 0, scale: 0.6 }}
								animate={{ opacity: 1, scale: 1 }}
								exit={{ opacity: 0, scale: 0.6 }}
								transition={SNAP_TRANSITION}
							>
								<Check className="size-3" />
							</motion.span>
						) : (
							<motion.span
								key="idle"
								initial={{ opacity: 0, scale: 0.6 }}
								animate={{ opacity: 1, scale: 1 }}
								exit={{ opacity: 0, scale: 0.6 }}
								transition={SNAP_TRANSITION}
							>
								<Copy className="size-3" />
							</motion.span>
						)}
					</AnimatePresence>
				</span>
				{copied ? "Copied" : "Copy"}
			</button>
		</div>
	);
}

export { InstallCommand };
