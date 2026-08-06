import { LASER_GRID, LASER_GRID_MASK, LASERS } from "#/constants";
import { cn } from "#/lib/utils.ts";
import type { Laser } from "#/types";

function laserAnimation(axis: Laser["axis"], dir: Laser["dir"]) {
	if (axis === "x") return dir === 1 ? "laser-x" : "laser-x-rev";
	return dir === 1 ? "laser-y" : "laser-y-rev";
}

function laserGradient(axis: Laser["axis"], dir: Laser["dir"]) {
	const head = "color-mix(in oklab, var(--color-primary) 55%, white)";
	const mid = "var(--color-primary)";
	if (axis === "x") {
		return dir === 1
			? `linear-gradient(90deg, transparent 0%, ${mid} 72%, ${head} 100%)`
			: `linear-gradient(270deg, transparent 0%, ${mid} 72%, ${head} 100%)`;
	}
	return dir === 1
		? `linear-gradient(180deg, transparent 0%, ${mid} 72%, ${head} 100%)`
		: `linear-gradient(0deg, transparent 0%, ${mid} 72%, ${head} 100%)`;
}

function LaserBeam({ laser }: { laser: Laser }) {
	const horizontal = laser.axis === "x";
	const anim = laserAnimation(laser.axis, laser.dir);

	return (
		<div
			className="absolute"
			style={
				horizontal
					? { top: laser.lane * LASER_GRID, left: 0, right: 0, height: 0 }
					: { left: laser.lane * LASER_GRID, top: 0, bottom: 0, width: 0 }
			}
		>
			{/* Full-lane traveler so translate % resolves against the panel, not the beam. */}
			<div
				className="laser-beam absolute inset-0"
				style={{
					["--laser-length" as string]: `${laser.length}px`,
					animationName: anim,
					animationDuration: `${laser.duration}s`,
					animationDelay: `${laser.phase}s`,
					animationTimingFunction: "linear",
					animationIterationCount: "infinite",
				}}
			>
				<span
					className="absolute top-0 left-0 block rounded-full opacity-80"
					style={{
						background: laserGradient(laser.axis, laser.dir),
						boxShadow:
							"0 0 6px color-mix(in oklab, var(--color-primary) 70%, transparent), 0 0 14px color-mix(in oklab, var(--color-primary) 35%, transparent)",
						...(horizontal
							? { width: laser.length, height: 2, marginTop: -1 }
							: { height: laser.length, width: 2, marginLeft: -1 }),
					}}
				/>
			</div>
		</div>
	);
}

// Grid motif with a handful of primary beams riding its lanes, faded out by a
// radial mask. Runs on CSS keyframes (see styles.css) so it costs nothing on
// the main thread; `prefers-reduced-motion` hides the beams there too.
function LaserGrid({ className }: { className?: string }) {
	return (
		<div
			className={cn("pointer-events-none absolute inset-0", className)}
			style={{
				maskImage: LASER_GRID_MASK,
				WebkitMaskImage: LASER_GRID_MASK,
			}}
			aria-hidden="true"
		>
			<div className="absolute inset-0 opacity-60 bg-[linear-gradient(var(--color-border)_1px,transparent_1px),linear-gradient(90deg,var(--color-border)_1px,transparent_1px)] bg-size-[24px_24px]" />
			<div className="absolute inset-0 overflow-hidden">
				{LASERS.map((laser) => (
					<LaserBeam
						key={`${laser.axis}-${laser.lane}-${laser.dir}`}
						laser={laser}
					/>
				))}
			</div>
		</div>
	);
}

export { LaserGrid };
