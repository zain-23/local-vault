/** Single laser riding a grid lane on a brand surface (auth panel, landing hero). */
export type Laser = {
	axis: "x" | "y";
	/** Grid line index (multiples of LASER_GRID px). */
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
