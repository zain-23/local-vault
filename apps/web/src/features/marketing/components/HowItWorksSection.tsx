import {
	ONBOARDING_STEPS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { Container } from "./Container.tsx";
import { Reveal } from "./Reveal.tsx";
import { Section } from "./Section.tsx";
import { SectionIntro } from "./SectionIntro.tsx";

function HowItWorksSection() {
	return (
		<Section id={SECTION_IDS.how}>
			<Container>
				<SectionIntro
					kicker="How it works"
					title="Three commands and your team is off Slack."
					description="No account to create, no card to enter, no server to run. The vault is a file on your machine and the key never leaves it."
				/>

				<div className="mt-11 grid grid-cols-1 gap-4.5 lg:grid-cols-3">
					{ONBOARDING_STEPS.map((step, index) => (
						<Reveal
							key={step.id}
							className="rounded-xl border border-border bg-card p-5.5"
						>
							<span className="inline-flex size-6.5 items-center justify-center rounded-md border border-accent-border bg-accent-soft font-mono text-xs font-medium text-primary">
								{index + 1}
							</span>
							<h3 className="mt-3.5 text-[15.5px] font-semibold tracking-tight">
								{step.title}
							</h3>
							<p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
								{step.description}
							</p>
							<code className="mt-3.5 inline-block rounded-sm border border-accent-border bg-accent-soft px-2 py-0.5 font-mono text-[12.5px] text-primary">
								{step.command}
							</code>
						</Reveal>
					))}
				</div>
			</Container>
		</Section>
	);
}

export { HowItWorksSection };
