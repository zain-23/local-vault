export type Laser = {
	axis: "x" | "y";
	lane: number;
	/** 1 = L→R / T→B, -1 = R→L / B→T */
	dir: 1 | -1;
	duration: number;
	phase: number;
	length: number;
};
