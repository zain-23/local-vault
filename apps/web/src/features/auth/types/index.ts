import type * as React from "react";

// One line of the brand terminal session: a typed command (`cmd`) or its
// output (`out`). `id` gives each line a stable React key.
export type TerminalEntry = { id: string } & (
	| { cmd: string }
	| { out: React.ReactNode }
);

/** Single laser riding a grid lane on the auth brand panel. */
export type Laser = {
	axis: "x" | "y";
	/** Grid line index (multiples of 24px). */
	lane: number;
	/** 1 = L→R / T→B, -1 = R→L / B→T */
	dir: 1 | -1;
	duration: number;
	/**
	 * Negative phase offset (seconds). Keeps stagger without a cold start —
	 * lasers are already mid-loop on first paint.
	 */
	phase: number;
	/** Trail length in px. */
	length: number;
};
