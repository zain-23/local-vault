import { motion } from "motion/react";
import { useState } from "react";

import type { ProductTabId } from "#/features/marketing/types";
import {
	PRODUCT_TABS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { SNAP_TRANSITION } from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";
import { ConsolePane } from "./ConsolePane.tsx";
import { Container } from "./Container.tsx";
import { ProductFrame } from "./ProductFrame.tsx";
import { Reveal } from "./Reveal.tsx";
import { Section } from "./Section.tsx";
import { SectionIntro } from "./SectionIntro.tsx";
import { TerminalPane } from "./TerminalPane.tsx";

/** Shared-layout id that slides the pill between tabs. */
const TAB_INDICATOR = "product-tab-indicator";

function ProductTourSection() {
	const [activeId, setActiveId] = useState<ProductTabId>("cli");
	const activeTab =
		PRODUCT_TABS.find((tab) => tab.id === activeId) ?? PRODUCT_TABS[0];

	return (
		<Section id={SECTION_IDS.product}>
			<Container>
				<SectionIntro
					kicker="Product tour"
					title="Two ways in — your terminal, or your browser."
					description={
						<>
							The <code className="text-primary">lv</code> CLI does the
							day-to-day work. The web console handles the things you'd rather
							click: approving devices, managing members, and reading the audit
							log.
						</>
					}
				/>

				<div className="mt-7.5">
					<Reveal className="mb-4.5 flex">
						<div
							role="tablist"
							aria-label="Product surfaces"
							className="inline-flex gap-0.5 rounded-lg border border-border bg-muted p-[3px]"
						>
							{PRODUCT_TABS.map((tab) => {
								const selected = tab.id === activeId;
								return (
									<button
										key={tab.id}
										id={`product-tab-${tab.id}`}
										type="button"
										role="tab"
										aria-selected={selected}
										aria-controls={`product-pane-${tab.id}`}
										onClick={() => setActiveId(tab.id)}
										className={cn(
											"relative z-[1] h-8 rounded-md px-4 text-[13.5px] font-medium transition-colors",
											selected
												? "text-foreground"
												: "text-muted-foreground hover:text-foreground",
										)}
									>
										{selected && (
											<motion.span
												layoutId={TAB_INDICATOR}
												transition={SNAP_TRANSITION}
												className="absolute inset-0 -z-10 rounded-md border border-border bg-card"
											/>
										)}
										{tab.label}
									</button>
								);
							})}
						</div>
					</Reveal>

					<ProductFrame tab={activeTab}>
						{activeId === "cli" ? <TerminalPane /> : <ConsolePane />}
					</ProductFrame>
				</div>
			</Container>
		</Section>
	);
}

export { ProductTourSection };
