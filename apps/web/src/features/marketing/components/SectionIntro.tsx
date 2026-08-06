import type * as React from "react";

import { Reveal } from "./Reveal.tsx";

// Kicker / headline / supporting line, each on its own staggered beat.
function SectionIntro({
	kicker,
	title,
	description,
}: {
	kicker: string;
	title: React.ReactNode;
	description: React.ReactNode;
}) {
	return (
		<>
			<Reveal>
				<span className="text-[13px] font-medium text-primary">{kicker}</span>
			</Reveal>
			<Reveal>
				<h2 className="text-[clamp(26px,3.4vw,38px)] leading-[1.12] font-semibold tracking-tight">
					{title}
				</h2>
			</Reveal>
			<Reveal>
				<p className="mt-3.5 max-w-[58ch] text-[16.5px] leading-relaxed text-muted-foreground">
					{description}
				</p>
			</Reveal>
		</>
	);
}

export { SectionIntro };
