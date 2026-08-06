import { motion } from "motion/react";
import type * as React from "react";

import {
	REVEAL_VIEWPORT,
	revealGroup,
} from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";

type SectionProps = React.ComponentProps<typeof motion.section>;

/**
 * A landing section: rule above, vertical rhythm, and the variant orchestrator
 * that staggers every `<Reveal>` inside it once the section scrolls into view.
 * `scroll-mt` keeps anchored jumps clear of the sticky nav.
 */
function Section({ className, ...props }: SectionProps) {
	return (
		<motion.section
			initial="hidden"
			whileInView="visible"
			viewport={REVEAL_VIEWPORT}
			variants={revealGroup}
			className={cn("scroll-mt-15 border-t border-border py-21", className)}
			{...props}
		/>
	);
}

export { Section };
