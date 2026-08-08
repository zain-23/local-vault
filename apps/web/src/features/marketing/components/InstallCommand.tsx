import { ArrowRight, Check, Copy, Download } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useId, useState } from "react";

import { useCopyToClipboard } from "#/features/marketing/hooks";
import type { InstallPlatformId } from "#/features/marketing/types";
import { INSTALL_PLATFORMS } from "#/features/marketing/utils/constants.ts";
import { SNAP_TRANSITION } from "#/features/marketing/utils/motion.ts";
import { cn } from "#/lib/utils.ts";

// Linux/macOS get the real `install.sh` one-liner with a self-confirming copy
// button; Windows has no installer script, so its tab links to GitHub instead.
function InstallCommand({ className }: { className?: string }) {
	const indicatorId = useId();
	const [platformId, setPlatformId] = useState<InstallPlatformId>(
		INSTALL_PLATFORMS[0].id,
	);
	const platform =
		INSTALL_PLATFORMS.find((p) => p.id === platformId) ?? INSTALL_PLATFORMS[0];

	return (
		<div
			className={cn(
				"inline-flex max-w-full flex-col items-center gap-2.5",
				className,
			)}
		>
			<div
				role="tablist"
				aria-label="Install platform"
				className="inline-flex gap-0.5 rounded-lg border border-border bg-muted p-[3px]"
			>
				{INSTALL_PLATFORMS.map((p) => {
					const selected = p.id === platformId;
					return (
						<button
							key={p.id}
							type="button"
							role="tab"
							aria-selected={selected}
							aria-controls={`install-command-${p.id}`}
							onClick={() => setPlatformId(p.id)}
							className={cn(
								"relative z-[1] h-7 rounded-md px-3 text-[12.5px] font-medium transition-colors",
								selected
									? "text-foreground"
									: "text-muted-foreground hover:text-foreground",
							)}
						>
							{selected && (
								<motion.span
									layoutId={`${indicatorId}-tab`}
									transition={SNAP_TRANSITION}
									className="absolute inset-0 -z-10 rounded-md border border-border bg-card"
								/>
							)}
							{p.label}
						</button>
					);
				})}
			</div>

			{platform.kind === "command" ? (
				<InstallCommandLine
					key={platform.id}
					id={`install-command-${platform.id}`}
					command={platform.command}
				/>
			) : (
				<InstallDownloadLink
					id={`install-command-${platform.id}`}
					platformLabel={platform.label}
					href={platform.href}
				/>
			)}
		</div>
	);
}

// Keyed by the caller so switching platforms remounts this — the copy
// button's "copied" state shouldn't survive onto a different command.
function InstallCommandLine({ id, command }: { id: string; command: string }) {
	const { copied, copy } = useCopyToClipboard(command);

	return (
		<div
			id={id}
			role="tabpanel"
			className="inline-flex max-w-full items-center gap-3 rounded-lg border border-border bg-card py-2 pr-2 pl-3.5 font-mono text-[13.5px]"
		>
			<span className="text-muted-foreground select-none">$</span>
			<span className="truncate text-foreground">{command}</span>
			<button
				type="button"
				onClick={copy}
				aria-label={`Copy "${command}" to the clipboard`}
				className={cn(
					"inline-flex h-7 flex-none items-center justify-center gap-1.5 rounded-md border border-border bg-secondary px-2.5 font-sans text-xs font-medium text-muted-foreground transition-colors",
					"hover:border-border-strong hover:text-foreground",
					copied && "border-success/40 text-success hover:text-success",
				)}
			>
				{/* Icons swap in place — same box, so the button never reflows. */}
				<span className="relative inline-flex size-3 items-center justify-center">
					<AnimatePresence initial={false} mode="popLayout">
						{copied ? (
							<motion.span
								key="done"
								initial={{ opacity: 0, scale: 0.6 }}
								animate={{ opacity: 1, scale: 1 }}
								exit={{ opacity: 0, scale: 0.6 }}
								transition={SNAP_TRANSITION}
							>
								<Check className="size-3" />
							</motion.span>
						) : (
							<motion.span
								key="idle"
								initial={{ opacity: 0, scale: 0.6 }}
								animate={{ opacity: 1, scale: 1 }}
								exit={{ opacity: 0, scale: 0.6 }}
								transition={SNAP_TRANSITION}
							>
								<Copy className="size-3" />
							</motion.span>
						)}
					</AnimatePresence>
				</span>
				{copied ? "Copied" : "Copy"}
			</button>
		</div>
	);
}

// No installer script for this platform (only release archives), so this
// links out to GitHub instead of showing a command that doesn't exist.
function InstallDownloadLink({
	id,
	platformLabel,
	href,
}: {
	id: string;
	platformLabel: string;
	href: string;
}) {
	return (
		<div id={id} role="tabpanel">
			<a
				href={href}
				target="_blank"
				rel="noreferrer noopener"
				className="inline-flex max-w-full items-center gap-2 rounded-lg border border-border bg-card py-2 pr-3 pl-3.5 text-[13.5px] font-medium text-foreground transition-colors hover:border-border-strong hover:bg-accent/40"
			>
				<Download className="size-3.5 text-muted-foreground" />
				Download for {platformLabel}
				<ArrowRight className="size-3.5 text-muted-foreground" />
			</a>
		</div>
	);
}

export { InstallCommand };
