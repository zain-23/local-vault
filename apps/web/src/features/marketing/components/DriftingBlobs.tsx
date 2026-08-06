import { motion } from "motion/react";
import type { CSSProperties } from "react";

import { cn } from "#/lib/utils.ts";

type Blob = {
	id: string;
	/** Size + placement. Inline style because these are one-off numbers. */
	style: CSSProperties;
	tint: string;
	/** Keyframe track: [from, to, back] on each axis. */
	drift: { x: number[]; y: number[]; scale: number[] };
	duration: number;
};

const tint = (token: string, percent: number) =>
	`radial-gradient(circle, color-mix(in oklab, var(${token}) ${percent}%, transparent), transparent 70%)`;

const GOLD: Blob = {
	id: "gold",
	style: { width: 640, height: 520, left: "50%", top: -260, marginLeft: -320 },
	tint: tint("--primary", 22),
	drift: { x: [0, -40, 0], y: [0, 26, 0], scale: [1, 1.08, 1] },
	duration: 26,
};

const AMBER: Blob = {
	id: "amber",
	style: { width: 480, height: 420, left: "12%", top: 60 },
	tint: tint("--warning", 13),
	drift: { x: [0, 56, 0], y: [0, -30, 0], scale: [1, 1.12, 1] },
	duration: 32,
};

const GREEN: Blob = {
	id: "green",
	style: { width: 420, height: 380, right: "8%", top: 10 },
	tint: tint("--success", 7),
	drift: { x: [0, -46, 0], y: [0, 34, 0], scale: [1, 0.92, 1] },
	duration: 38,
};

export const HERO_BLOBS: Blob[] = [GOLD, AMBER, GREEN];

/** Same wash, pulled further off the top edge so the CTA copy stays readable. */
export const CTA_BLOBS: Blob[] = [
	{ ...GOLD, style: { ...GOLD.style, top: -320 } },
	{ ...AMBER, style: { ...AMBER.style, left: "20%", top: -40 } },
];

/**
 * Slow-drifting colour wash behind the hero and the closing CTA. Only
 * `transform` animates, so it stays on the compositor; the page-level
 * `MotionConfig reducedMotion="user"` freezes it for anyone who asked for less.
 */
function DriftingBlobs({
	blobs = HERO_BLOBS,
	className,
}: {
	blobs?: Blob[];
	className?: string;
}) {
	return (
		<div
			className={cn(
				"pointer-events-none absolute inset-0 overflow-hidden",
				className,
			)}
			aria-hidden="true"
		>
			{blobs.map((blob) => (
				<motion.div
					key={blob.id}
					className="absolute rounded-full blur-[70px] will-change-transform"
					style={{ ...blob.style, background: blob.tint }}
					animate={blob.drift}
					transition={{
						duration: blob.duration,
						repeat: Number.POSITIVE_INFINITY,
						ease: "easeInOut",
					}}
				/>
			))}
		</div>
	);
}

export { DriftingBlobs, type Blob };
