/**
 * One fractal-noise filter for the whole page. Mounted once by `LandingPage`;
 * every `<GrainOverlay>` just references it by id, so the (expensive) turbulence
 * is only ever rasterised once.
 */
const GRAIN_FILTER_ID = "lv-grain";

function GrainFilterDefs() {
	return (
		<svg className="absolute size-0" aria-hidden="true">
			<filter id={GRAIN_FILTER_ID}>
				<feTurbulence
					type="fractalNoise"
					baseFrequency="0.82"
					numOctaves="3"
					stitchTiles="stitch"
				/>
				<feColorMatrix type="saturate" values="0" />
			</filter>
		</svg>
	);
}

/** Film grain over a section's gradients. Needs a positioned ancestor. */
function GrainOverlay() {
	return (
		<svg
			className="pointer-events-none absolute inset-0 z-[2] opacity-[0.055] mix-blend-overlay"
			width="100%"
			height="100%"
			aria-hidden="true"
		>
			<rect width="100%" height="100%" filter={`url(#${GRAIN_FILTER_ID})`} />
		</svg>
	);
}

export { GrainFilterDefs, GrainOverlay };
