import { motion } from "motion/react";

/** Laptop → relay, then relay → the other laptop. */
const HOP_UP = "M74 64 C140 64 160 28 204 22";
const HOP_DOWN = "M256 22 C300 28 320 64 386 64";
/** Seconds for one packet to cross a hop. */
const HOP_DURATION = 2.8;

// A short lit segment sliding along a hop — the sealed blob in transit.
function Packet({ d, delay }: { d: string; delay: number }) {
	return (
		<motion.path
			d={d}
			fill="none"
			strokeWidth={2}
			strokeLinecap="round"
			className="stroke-primary"
			initial={{ pathLength: 0.14, pathOffset: 0 }}
			animate={{ pathLength: 0.14, pathOffset: 1 }}
			transition={{
				duration: HOP_DURATION,
				repeat: Number.POSITIVE_INFINITY,
				ease: "linear",
				delay,
			}}
		/>
	);
}

function DeviceGlyph({ x, label }: { x: number; label: string }) {
	return (
		<>
			<rect
				x={x}
				y={52}
				width={66}
				height={34}
				rx={4}
				strokeWidth={1.2}
				fill="none"
				className="stroke-border-strong"
			/>
			<rect
				x={x + 11}
				y={60}
				width={44}
				height={18}
				rx={2}
				strokeWidth={1}
				fill="none"
				className="stroke-border"
			/>
			<text
				x={x + 33}
				y={96}
				textAnchor="middle"
				fontSize={8.5}
				className="fill-muted-foreground font-mono"
			>
				{label}
			</text>
		</>
	);
}

/**
 * Two paired devices exchanging a sealed blob through the relay. The relay is
 * drawn dashed on purpose: it's the part of the system you don't have to trust.
 */
function RelayDiagram() {
	return (
		<svg
			viewBox="0 0 460 96"
			fill="none"
			className="block h-auto w-full"
			role="img"
			aria-label="Two paired devices exchanging an encrypted blob through a relay that only ever sees ciphertext"
		>
			<DeviceGlyph x={8} label="your laptop" />
			<DeviceGlyph x={386} label="dana's mac" />

			<rect
				x={204}
				y={10}
				width={52}
				height={26}
				rx={4}
				strokeWidth={1.2}
				strokeDasharray="3 3"
				className="stroke-border-strong"
			/>
			<text
				x={230}
				y={27}
				textAnchor="middle"
				fontSize={8.5}
				className="fill-muted-foreground font-mono"
			>
				relay
			</text>

			<path d={HOP_UP} strokeWidth={1.2} className="stroke-border" />
			<path d={HOP_DOWN} strokeWidth={1.2} className="stroke-border" />

			<text
				x={230}
				y={52}
				textAnchor="middle"
				fontSize={7.5}
				letterSpacing={1.2}
				className="fill-muted-foreground/60 font-mono"
			>
				CIPHERTEXT
			</text>

			<Packet d={HOP_UP} delay={0} />
			<Packet d={HOP_DOWN} delay={HOP_DURATION / 2} />
		</svg>
	);
}

export { RelayDiagram };
