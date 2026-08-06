import { motion } from "motion/react";
import type * as React from "react";

import { revealItem } from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";

type RevealProps = React.ComponentProps<typeof motion.div>;

/**
 * A single staggered reveal. Deliberately has no `initial`/`animate` of its own
 * so it inherits the hidden→visible label from the nearest `<Section>`, which is
 * what drives the stagger. Drop it anywhere inside one, at any nesting depth.
 */
function Reveal({ className, ...props }: RevealProps) {
	return (
		<motion.div variants={revealItem} className={cn(className)} {...props} />
	);
}

export { Reveal };
