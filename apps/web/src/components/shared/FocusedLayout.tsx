import type * as React from "react";

import { LVLogo } from "#/components/shared";

const FOOTER_LINKS = ["Privacy", "Terms", "Security"] as const;

// Centered page chrome for focused flows outside the app shell (CLI device
// link, workspace invite accept, etc.): logo header, grid backdrop, footer.
function FocusedLayout({ children }: { children: React.ReactNode }) {
	return (
		<div className="relative flex min-h-svh flex-col bg-background text-foreground">
			{/* Grid backdrop, matched to the auth brand panel so the flow feels of a
          piece with sign-in. Masked to a soft glow anchored top-left. */}
			<div
				className="pointer-events-none absolute inset-0 opacity-60 bg-[linear-gradient(var(--color-border)_1px,transparent_1px),linear-gradient(90deg,var(--color-border)_1px,transparent_1px)] bg-size-[24px_24px]"
				style={{
					maskImage:
						"radial-gradient(ellipse at 25% 15%, black 0%, transparent 72%)",
					WebkitMaskImage:
						"radial-gradient(ellipse at 25% 15%, black 0%, transparent 72%)",
				}}
				aria-hidden="true"
			/>

			<header className="relative flex items-center justify-between px-6 py-5 sm:px-10">
				<LVLogo size={35} />
				<button
					type="button"
					className="text-[13px] text-muted-foreground transition-colors hover:text-foreground"
				>
					Need help?
				</button>
			</header>

			<main className="relative flex flex-1 items-center justify-center px-6 py-10">
				<div className="w-full max-w-115">{children}</div>
			</main>

			<footer className="relative flex flex-wrap items-center justify-center gap-x-4 gap-y-1 px-6 py-6 text-[11.5px] text-muted-foreground">
				<span className="font-mono">© 2026 LocalVault</span>
				{FOOTER_LINKS.map((label) => (
					<button
						key={label}
						type="button"
						className="transition-colors hover:text-foreground"
					>
						{label}
					</button>
				))}
			</footer>
		</div>
	);
}

export { FocusedLayout };
