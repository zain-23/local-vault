import type * as React from "react";

import { cn } from "#/lib/utils.ts";

// Rounded icon badge that anchors the confirmation-style screens (magic link,
// check inbox, verify email, 2FA, success). `tone` tints the icon.
function AuthIconBadge({
	children,
	tone = "accent",
	className,
}: {
	children: React.ReactNode;
	tone?: "accent" | "success";
	className?: string;
}) {
	return (
		<div
			className={cn(
				"mb-4.5 flex size-12 items-center justify-center rounded-2xl border border-border bg-muted",
				tone === "success" ? "text-success" : "text-primary",
				className,
			)}
		>
			{children}
		</div>
	);
}

export { AuthIconBadge };
