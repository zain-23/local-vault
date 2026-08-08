import { ArrowRight, Check, Info } from "lucide-react";
import { motion } from "motion/react";

import { LaserGrid } from "#/components/shared/LaserGrid.tsx";
import { Button } from "#/components/ui/Button.tsx";
import {
	HERO_TRUST_POINTS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { revealGroup } from "#/features/marketing/utils/motion.ts";
import { Container } from "./Container.tsx";
import { InstallCommand } from "./InstallCommand.tsx";
import { Reveal } from "./Reveal.tsx";

function HeroSection() {
	return (
		<section className="relative overflow-hidden pt-19 pb-24">
			{/* <DriftingBlobs /> */}
			<LaserGrid className="z-[1]" />
			{/* <GrainOverlay /> */}

			{/* Plays on load rather than on scroll — it's already in view. */}
			<motion.div
				initial="hidden"
				animate="visible"
				variants={revealGroup}
				className="relative z-[3]"
			>
				<Container className="text-center">
					<Reveal>
						<h1 className="mx-auto mt-5 max-w-[19ch] text-[clamp(34px,5.6vw,62px)] leading-[1.04] font-semibold tracking-[-0.035em]">
							Stop sharing secrets{" "}
							<span
								className="bg-clip-text text-transparent"
								style={{ backgroundImage: "linear-gradient(180deg, var(--foreground) 30%, color-mix(in oklab, var(--primary) 72%, #000))" }}
							>
								over Slack.
							</span>
						</h1>
					</Reveal>

					<Reveal>
						<p className="mx-auto mt-5.5 max-w-[54ch] text-[17.5px] leading-relaxed text-muted-foreground">
							LocalVault replaces{" "}
							<b className="font-medium text-foreground">.env</b> files with an
							encrypted vault that syncs peer-to-peer between your team's
							machines. Your secrets are sealed before they leave, so there's
							nothing in the middle worth stealing.
						</p>
					</Reveal>

					<Reveal className="mt-7.5 flex flex-wrap items-center justify-center gap-2.5">
						<Button asChild size="lg">
							<a href={`#${SECTION_IDS.install}`}>
								Get started
								<ArrowRight />
							</a>
						</Button>
						<Button asChild size="lg" variant="outline">
							<a href={`#${SECTION_IDS.how}`}>
								<Info />
								How it works
							</a>
						</Button>
					</Reveal>

					<Reveal>
						<InstallCommand className="mt-6.5" />
					</Reveal>

					<Reveal className="mt-6.5 flex flex-wrap items-center justify-center gap-x-5.5 gap-y-2 text-[13px] text-muted-foreground">
						{HERO_TRUST_POINTS.map((point) => (
							<span key={point.id} className="inline-flex items-center gap-1.5">
								<Check className="size-3.5 text-success" />
								{point.label}
							</span>
						))}
					</Reveal>
				</Container>
			</motion.div>
		</section>
	);
}

export { HeroSection };
