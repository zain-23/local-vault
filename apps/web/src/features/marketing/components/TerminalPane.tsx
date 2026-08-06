import { motion } from "motion/react";

import type { TerminalTone } from "#/features/marketing/types";
import { TERMINAL_TRANSCRIPT } from "#/features/marketing/utils/constants.ts";
import {
	terminalGroup,
	terminalLine,
} from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";

const TONE_CLASS: Record<TerminalTone, string> = {
	prompt: "text-primary",
	success: "text-success",
	warning: "text-warning",
	muted: "text-muted-foreground",
};

// Block cursor. Steps rather than a fade, so it reads as a real terminal.
function Caret() {
	return (
		<motion.span
			className="ml-1 inline-block h-3.5 w-[7px] -mb-0.5 bg-primary align-text-bottom"
			animate={{ opacity: [1, 1, 0, 0] }}
			transition={{
				duration: 1.05,
				times: [0, 0.499, 0.5, 1],
				repeat: Number.POSITIVE_INFINITY,
				ease: "linear",
			}}
			aria-hidden="true"
		/>
	);
}

/**
 * Replays the transcript line by line every time it mounts — which is exactly
 * when you want it: first scroll into view, and again on each return to the tab.
 */
function TerminalPane() {
	return (
		<motion.div
			initial="hidden"
			whileInView="visible"
			viewport={{ once: true, amount: 0.2 }}
			variants={terminalGroup}
			className="min-h-[330px] px-5.5 pt-5 pb-6.5 font-mono text-[13px] leading-[1.95]"
		>
			{TERMINAL_TRANSCRIPT.map((line) => (
				<motion.div key={line.id} variants={terminalLine}>
					{line.segments.map((segment, index) => (
						<span
							key={`${line.id}-${segment.text}`}
							className={cn(
								segment.tone && TONE_CLASS[segment.tone],
								index > 0 && "ml-[0.5ch]",
							)}
						>
							{segment.text}
						</span>
					))}
					{line.caret ? <Caret /> : null}
				</motion.div>
			))}
		</motion.div>
	);
}

export { TerminalPane };
