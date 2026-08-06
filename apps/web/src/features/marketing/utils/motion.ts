import {
	stagger,
	type Transition,
	type Variants,
	type ViewportOptions,
} from "motion/react";

/**
 * Springs, not eases. Low bounce throughout — this is security software, it
 * shouldn't wobble. `visualDuration` is the time to *visually* settle, so these
 * read the same as the durations they replace.
 */
export const REVEAL_TRANSITION: Transition = {
	type: "spring",
	visualDuration: 0.8,
	bounce: 0,
};

/** Snappier, with a touch of overshoot — used for the tab indicator. */
export const SNAP_TRANSITION: Transition = {
	type: "spring",
	visualDuration: 0.5,
	bounce: 0.18,
};

/** Reveal once, a little before the element is fully on screen. */
export const REVEAL_VIEWPORT: ViewportOptions = {
	once: true,
	amount: 0.15,
	margin: "0px 0px -5% 0px",
};

/** Gap between sibling reveals inside one section (seconds). */
const REVEAL_STAGGER = 0.06;

/**
 * Orchestrator. Put it on the section; every descendant `motion` element that
 * doesn't declare its own `animate` inherits the label and staggers in.
 */
export const revealGroup: Variants = {
	hidden: {},
	visible: { transition: { delayChildren: stagger(REVEAL_STAGGER) } },
};

/** The thing that actually moves. Pairs with `revealGroup`. */
export const revealItem: Variants = {
	hidden: { opacity: 0, y: 18 },
	visible: { opacity: 1, y: 0, transition: REVEAL_TRANSITION },
};

/** Product frame tilting upright as it enters the viewport. */
export const productFrame: Variants = {
	hidden: { opacity: 0, rotateX: 7, y: 34, scale: 0.975 },
	visible: {
		opacity: 1,
		rotateX: 0,
		y: 0,
		scale: 1,
		transition: { type: "spring", visualDuration: 1, bounce: 0 },
	},
};

/** Terminal transcript: lines land one after another, like they're being typed. */
export const terminalGroup: Variants = {
	hidden: {},
	visible: {
		transition: { delayChildren: stagger(0.3, { startDelay: 0.26 }) },
	},
};

export const terminalLine: Variants = {
	hidden: { opacity: 0, y: 3 },
	visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: "easeOut" } },
};

/** Console table rows, tighter than the transcript. */
export const tableGroup: Variants = {
	hidden: {},
	visible: { transition: { delayChildren: stagger(0.07) } },
};

export const tableRow: Variants = {
	hidden: { opacity: 0, y: 6 },
	visible: { opacity: 1, y: 0, transition: { duration: 0.5, ease: "easeOut" } },
};
