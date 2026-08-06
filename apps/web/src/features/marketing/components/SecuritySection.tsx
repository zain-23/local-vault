import {
	CRYPTO_PRIMITIVES,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { Container } from "./Container.tsx";
import { Reveal } from "./Reveal.tsx";
import { Section } from "./Section.tsx";
import { SectionIntro } from "./SectionIntro.tsx";

function SecuritySection() {
	return (
		<Section id={SECTION_IDS.security}>
			<Container>
				<SectionIntro
					kicker="Security"
					title="Boring cryptography, on purpose."
					description="Nothing here is homegrown. Every primitive is a standard you can look up, and every key it produces stays on the machine that made it."
				/>

				<div className="mt-10 grid grid-cols-1 gap-4.5 sm:grid-cols-2 lg:grid-cols-4">
					{CRYPTO_PRIMITIVES.map((primitive) => (
						<Reveal
							key={primitive.id}
							className="rounded-xl border border-border bg-card p-4.5"
						>
							<div className="font-mono text-[13.5px] text-primary">
								{primitive.name}
							</div>
							<div className="mt-2.5 text-[11px] tracking-[0.07em] text-muted-foreground uppercase">
								{primitive.role}
							</div>
							<p className="mt-2 text-[13.5px] leading-snug text-muted-foreground">
								{primitive.description}
							</p>
						</Reveal>
					))}
				</div>
			</Container>
		</Section>
	);
}

export { SecuritySection };
