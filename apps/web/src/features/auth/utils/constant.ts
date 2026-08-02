import type { Laser } from "../types";

/** Client-side cooldown before another password-reset / magic-link email can be sent. */
export const RESEND_COOLDOWN_SECONDS = 60;

/** Auth brand panel grid cell size (px). */
export const LASER_GRID = 24;

/** Soft radial fade shared by the grid motif and lasers. */
export const LASER_GRID_MASK =
	"radial-gradient(ellipse at 25% 15%, black 0%, transparent 72%)";

/**
 * Sparse primary lasers, snapped to the 24px grid.
 * 2 per edge (L/R/T/B), lanes spread evenly so no corner clusters.
 */
export const LASERS: Laser[] = [
	// ← from left
	{ axis: "x", lane: 4, dir: 1, duration: 9, phase: -2.1, length: 48 },
	{ axis: "x", lane: 18, dir: 1, duration: 11, phase: -6.4, length: 40 },
	// → from right
	{ axis: "x", lane: 10, dir: -1, duration: 10, phase: -3.8, length: 44 },
	{ axis: "x", lane: 24, dir: -1, duration: 12, phase: -7.5, length: 36 },
	// ↓ from top
	{ axis: "y", lane: 4, dir: 1, duration: 10.5, phase: -1.5, length: 40 },
	{ axis: "y", lane: 20, dir: 1, duration: 9.5, phase: -5.2, length: 48 },
	// ↑ from bottom
	{ axis: "y", lane: 12, dir: -1, duration: 11.5, phase: -4.0, length: 36 },
	{ axis: "y", lane: 28, dir: -1, duration: 10, phase: -8.1, length: 44 },
];
