import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useMotionValueEvent, useScroll } from "motion/react";
import { useState } from "react";

import { LVLogo } from "#/components/shared/LVLogo.tsx";
import { Button } from "#/components/ui/Button.tsx";
import { meQuery } from "#/features/auth/api";
import {
	NAV_LINKS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { cn } from "#/lib/utils.ts";
import { Container } from "./Container.tsx";

/** Scroll distance past which the bar earns its background. */
const STUCK_AFTER_PX = 6;

function LandingNav() {
	const { scrollY } = useScroll();
	const [stuck, setStuck] = useState(false);
	const { data: user } = useQuery(meQuery);

	// setState only on the crossing, not on every scroll frame.
	useMotionValueEvent(scrollY, "change", (y) => {
		const next = y > STUCK_AFTER_PX;
		setStuck((prev) => (prev === next ? prev : next));
	});

	return (
		<header
			className={cn(
				"sticky top-0 z-50 border-b border-transparent transition-colors duration-300",
				stuck && "border-border bg-background/70 backdrop-blur-md",
			)}
		>
			<Container className="flex h-15 items-center justify-between gap-6">
				<a href={`#${SECTION_IDS.top}`} aria-label="LocalVault home">
					<LVLogo size={24} className="gap-2.5" />
				</a>

				<nav className="hidden items-center gap-1 lg:flex">
					{NAV_LINKS.map((link) => (
						<a
							key={link.href + link.label}
							href={link.href}
							className="rounded-md px-2.5 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
						>
							{link.label}
						</a>
					))}
				</nav>

				<div className="flex items-center gap-2">
					{user ? (
						<Button asChild>
							<Link to="/dashboard">Dashboard</Link>
						</Button>
					) : (
						<>
							<Button asChild variant="ghost" className="hidden sm:inline-flex">
								<Link to="/auth/login">Sign in</Link>
							</Button>
							<Button asChild>
								<a href={`#${SECTION_IDS.install}`}>Get started</a>
							</Button>
						</>
					)}
				</div>
			</Container>
		</header>
	);
}

export { LandingNav };
