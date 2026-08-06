import { MotionConfig } from "motion/react";

import { SECTION_IDS } from "#/features/marketing/utils/constants.ts";
import { FeaturesSection } from "./FeaturesSection.tsx";
import { GrainFilterDefs } from "./Grain.tsx";
import { HeroSection } from "./HeroSection.tsx";
import { HowItWorksSection } from "./HowItWorksSection.tsx";
import { InstallCtaSection } from "./InstallCtaSection.tsx";
import { LandingFooter } from "./LandingFooter.tsx";
import { LandingNav } from "./LandingNav.tsx";
import { ProductTourSection } from "./ProductTourSection.tsx";
import { SecuritySection } from "./SecuritySection.tsx";

/**
 * Public marketing page at `/`.
 *
 * `reducedMotion="user"` is set once here rather than branching in every
 * component: Motion then drops transform and layout animation for visitors who
 * asked for less motion, and everything below can be written as if it's always on.
 */
function LandingPage() {
	return (
		<MotionConfig reducedMotion="user">
			<GrainFilterDefs />
			<LandingNav />
			<main id={SECTION_IDS.top}>
				<HeroSection />
				<ProductTourSection />
				<HowItWorksSection />
				<FeaturesSection />
				<SecuritySection />
				<InstallCtaSection />
			</main>
			<LandingFooter />
		</MotionConfig>
	);
}

export { LandingPage };
