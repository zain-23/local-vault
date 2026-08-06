import { ArrowRight, Check } from "lucide-react";

import { GithubIcon } from "#/components/shared/GithubIcon.tsx";
import { Badge } from "#/components/ui/Badge.tsx";
import { Button } from "#/components/ui/Button.tsx";
import {
	EXTERNAL_LINKS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { Container } from "./Container.tsx";
import { CTA_BLOBS, DriftingBlobs } from "./DriftingBlobs.tsx";
import { GrainOverlay } from "./Grain.tsx";
import { InstallCommand } from "./InstallCommand.tsx";
import { Reveal } from "./Reveal.tsx";
import { Section } from "./Section.tsx";

function InstallCtaSection() {
	return (
		<Section
			id={SECTION_IDS.install}
			className="relative overflow-hidden py-24 text-center"
		>
			<DriftingBlobs blobs={CTA_BLOBS} />
			<GrainOverlay />

			<Container className="relative z-[3]">
				<Reveal className="inline-block">
					<Badge className="rounded-md border-accent-border bg-accent-soft text-primary">
						<Check className="size-3" />
						Free and open source
					</Badge>
				</Reveal>

				<Reveal>
					<h2 className="mx-auto mt-4 max-w-[22ch] text-[clamp(28px,4.2vw,46px)] leading-[1.08] font-semibold tracking-[-0.03em]">
						Your next secret shouldn't touch a chat window.
					</h2>
				</Reveal>

				<Reveal>
					<p className="mx-auto mt-4.5 max-w-[48ch] text-[17px] text-muted-foreground">
						One install, one vault, every device you actually trust. Takes about
						two minutes.
					</p>
				</Reveal>

				<Reveal className="mt-7.5 flex flex-wrap items-center justify-center gap-2.5">
					<Button asChild size="lg">
						<a
							href={EXTERNAL_LINKS.releases}
							target="_blank"
							rel="noreferrer noopener"
						>
							Install LocalVault
							<ArrowRight />
						</a>
					</Button>
					<Button asChild size="lg" variant="outline">
						<a
							href={EXTERNAL_LINKS.github}
							target="_blank"
							rel="noreferrer noopener"
						>
							<GithubIcon />
							Star on GitHub
						</a>
					</Button>
				</Reveal>

				<Reveal>
					<InstallCommand className="mt-6" />
				</Reveal>
			</Container>
		</Section>
	);
}

export { InstallCtaSection };
