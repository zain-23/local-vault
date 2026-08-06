import { AnimatePresence, motion } from "motion/react";
import type * as React from "react";

import { Badge } from "#/components/ui/Badge.tsx";
import type { ProductTab } from "#/features/marketing/types";
import { productFrame } from "#/features/marketing/utils/motion.ts";

/** Window chrome dots — flat, not the macOS traffic lights. */
function ChromeDots() {
	return (
		<>
			<span className="size-2.5 rounded-full bg-border-strong" />
			<span className="size-2.5 rounded-full bg-border-strong" />
			<span className="size-2.5 rounded-full bg-border-strong" />
		</>
	);
}

/**
 * The app window the product tour renders inside. Tilts upright from a slight
 * perspective rotation the first time it scrolls into view; the pane inside then
 * crossfades whenever the active tab changes.
 */
function ProductFrame({
	tab,
	children,
}: {
	tab: ProductTab;
	children: React.ReactNode;
}) {
	return (
		<motion.div
			initial="hidden"
			whileInView="visible"
			viewport={{ once: true, amount: 0.2 }}
			variants={productFrame}
			style={{ transformPerspective: 1400, transformOrigin: "50% 0" }}
			className="overflow-hidden rounded-xl border border-border bg-card shadow-[0_40px_80px_-30px_rgba(0,0,0,0.7)]"
		>
			<div className="flex items-center gap-2 border-b border-border bg-sidebar px-3.5 py-2.5">
				<ChromeDots />
				<span className="ml-1.5 font-mono text-[11.5px] text-muted-foreground">
					{tab.frameTitle}
				</span>
				<Badge className="ml-auto rounded-md border-success/35 bg-success/10 text-success">
					{tab.frameBadge}
				</Badge>
			</div>

			<AnimatePresence mode="wait" initial={false}>
				<motion.div
					key={tab.id}
					id={`product-pane-${tab.id}`}
					role="tabpanel"
					aria-labelledby={`product-tab-${tab.id}`}
					initial={{ opacity: 0 }}
					animate={{ opacity: 1 }}
					exit={{ opacity: 0 }}
					transition={{ duration: 0.18, ease: "easeOut" }}
				>
					{children}
				</motion.div>
			</AnimatePresence>
		</motion.div>
	);
}

export { ProductFrame };
