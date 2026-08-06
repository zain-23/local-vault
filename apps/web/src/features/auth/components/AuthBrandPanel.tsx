import { Github, Link2, Lock } from "lucide-react";

import { LaserGrid } from "#/components/shared/LaserGrid.tsx";
import { LVLogo } from "#/components/shared/LVLogo.tsx";
import { BrandTerminal } from "./BrandTerminal.tsx";

// Trust chips shown under the terminal signature.
const CHIPS = [
	{ label: "Zero-knowledge", icon: Lock },
	{ label: "Peer-to-peer", icon: Link2 },
	{ label: "MIT licensed", icon: Github },
];

function AuthBrandPanel() {
	return (
		<div className="relative flex h-full flex-col overflow-hidden bg-muted px-13 py-12">
			<LaserGrid />
			{/* soft gold glow */}
			<div
				className="pointer-events-none absolute -top-30 -left-15 size-90 rounded-full bg-accent-soft blur-[70px]"
				aria-hidden="true"
			/>

			<div className="relative">
				<LVLogo size={40} />
			</div>

			<div className="relative flex flex-1 flex-col justify-center">
				<h2 className="text-6xl font-semibold">
					Secrets that never
					<br />
					leave your machine.
				</h2>
				<p className="text-muted-foreground mt-4">
					Encrypt on-device, sync peer-to-peer. We relay ciphertext and nothing
					else.
				</p>

				<BrandTerminal className="mt-10" />

				<div className="mt-5.5 flex flex-wrap gap-2">
					{CHIPS.map(({ label, icon: Icon }) => (
						<span
							key={label}
							className="inline-flex items-center gap-1.5 rounded-full border border-border bg-background px-2.5 py-1 text-[11.5px] text-muted-foreground"
						>
							<Icon className="size-3 text-primary" />
							{label}
						</span>
					))}
				</div>
			</div>

			<div className="relative font-mono text-[11.5px] text-muted-foreground/70">
				© 2026 LocalVault
			</div>
		</div>
	);
}

export { AuthBrandPanel };
