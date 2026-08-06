import type * as React from "react";
import { AuthBrandPanel } from "./AuthBrandPanel.tsx";

const FOOTER_LINKS = ["Privacy", "Terms", "Security"] as const;

function AuthLayout({ children }: { children: React.ReactNode }) {
	return (
		<div className="grid min-h-svh bg-background text-foreground lg:grid-cols-[1.1fr_1fr]">
			<div className="hidden lg:block">
				<AuthBrandPanel />
			</div>

			<div className="relative flex flex-col overflow-hidden px-6 py-8 sm:px-12">
				<header className="relative flex items-center justify-end">
					<button
						type="button"
						className="text-[13px] text-muted-foreground transition-colors hover:text-foreground"
					>
						Need help?
					</button>
				</header>

				<main className="relative flex flex-1 items-center justify-center py-10">
					<div className="w-full max-w-90">{children}</div>
				</main>

				<footer className="relative flex flex-wrap items-center justify-center gap-x-4 gap-y-1 text-[11.5px] text-muted-foreground">
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
		</div>
	);
}

export { AuthLayout };
