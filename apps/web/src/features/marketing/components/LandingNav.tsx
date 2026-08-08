import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { Star } from "lucide-react";
import { useMotionValueEvent, useScroll } from "motion/react";
import { useState } from "react";

import { GithubIcon } from "#/components/shared/GithubIcon.tsx";
import { LVLogo } from "#/components/shared/LVLogo.tsx";
import { Button } from "#/components/ui/Button.tsx";
import { meQuery } from "#/features/auth/api";
import { useGithubStars } from "#/features/marketing/hooks";
import {
	EXTERNAL_LINKS,
	NAV_LINKS,
	SECTION_IDS,
} from "#/features/marketing/utils/constants.ts";
import { cn } from "#/lib/utils.ts";
import { Container } from "./Container.tsx";

const starFormatter = new Intl.NumberFormat("en-US", { notation: "compact" });

/** Scroll distance past which the bar earns its background. */
const STUCK_AFTER_PX = 6;

function LandingNav() {
	const { scrollY } = useScroll();
	const [stuck, setStuck] = useState(false);
	const { data: user } = useQuery(meQuery);
	const { data: starCount } = useGithubStars();

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
							target={link.external ? "_blank" : undefined}
							rel={link.external ? "noreferrer noopener" : undefined}
							className="rounded-md px-2.5 py-1.5 text-sm text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
						>
							{link.label}
						</a>
					))}
				</nav>

				<div className="flex items-center gap-2">
					<Button asChild variant="outline" className="hidden sm:inline-flex">
						<a
							href={EXTERNAL_LINKS.github}
							target="_blank"
							rel="noreferrer noopener"
							aria-label="Star LocalVault on GitHub"
						>
							<GithubIcon size={14} />
							<Star className="size-3.5 fill-current" />
							{starCount !== undefined ? starFormatter.format(starCount) : ""}
						</a>
					</Button>

					{user ? (
						<Button asChild>
							<Link to="/dashboard">Dashboard</Link>
						</Button>
					) : (
						<Button asChild variant="ghost">
							<Link to="/auth/login">Sign in</Link>
						</Button>
					)}
				</div>
			</Container>
		</header>
	);
}

export { LandingNav };
