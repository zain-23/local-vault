import type * as React from "react";

import type { Feature } from "#/features/marketing/types";
import {
	FEATURES,
	RELAY_FEATURE,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { cn } from "#/lib/utils.ts";
import { Container } from "./Container.tsx";
import { RelayDiagram } from "./RelayDiagram.tsx";
import { RelayLedger } from "./RelayLedger.tsx";
import { Reveal } from "./Reveal.tsx";
import { Section } from "./Section.tsx";
import { SectionIntro } from "./SectionIntro.tsx";

/** Splits `$ lv pair` into a dim `$ ` prefix and a highlighted command. */
function FeatureCommand({ feature }: { feature: Feature }) {
	if (!feature.command) return null;

	const accent = feature.commandAccent;
	const at = accent ? feature.command.indexOf(accent) : -1;
	const before = at < 0 ? feature.command : feature.command.slice(0, at);
	const after =
		at < 0 || !accent ? "" : feature.command.slice(at + accent.length);

	return (
		<div className="mt-3.5 border-t border-border pt-3 font-mono text-xs text-muted-foreground">
			{before}
			{at >= 0 ? <span className="text-primary">{accent}</span> : null}
			{after}
		</div>
	);
}

function FeatureCard({
	feature,
	className,
	children,
}: {
	feature: Feature;
	className?: string;
	children?: React.ReactNode;
}) {
	const { icon: Icon } = feature;

	return (
		<Reveal
			className={cn(
				"relative overflow-hidden rounded-xl border border-border bg-card p-5.5",
				className,
			)}
		>
			<span className="inline-flex size-8.5 items-center justify-center rounded-lg border border-accent-border bg-accent-soft text-primary">
				<Icon className="size-[17px]" />
			</span>
			<h3 className="mt-3.5 text-[15.5px] font-semibold tracking-tight">
				{feature.title}
			</h3>
			<p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
				{feature.description}
			</p>
			{children}
			<FeatureCommand feature={feature} />
		</Reveal>
	);
}

function FeaturesSection() {
	return (
		<Section id={SECTION_IDS.features} className="relative overflow-hidden">
			<Container>
				<SectionIntro
					kicker="Features"
					title={
						<>
							Everything a shared <span className="text-primary">.env</span>{" "}
							should have been.
						</>
					}
					description="Built for small teams who need to move fast without leaving credentials scattered across a dozen chat threads."
				/>

				<div className="mt-11 grid grid-cols-1 gap-4.5 sm:grid-cols-2 lg:grid-cols-3">
					<FeatureCard feature={RELAY_FEATURE} className="sm:col-span-2">
						<div className="mt-4">
							<RelayDiagram />
						</div>
						<RelayLedger />
					</FeatureCard>

					{FEATURES.map((feature) => (
						<FeatureCard key={feature.id} feature={feature} />
					))}
				</div>
			</Container>
		</Section>
	);
}

export { FeaturesSection };
