import type * as React from "react";

import { cn } from "#/lib/utils.ts";

function LVLogo({
	size = 22,
	withWord = true,
	className,
	...props
}: React.ComponentProps<"span"> & {
	size?: number;
	withWord?: boolean;
}) {
	return (
		<span
			className={cn(
				"inline-flex items-center gap-2 text-foreground",
				className,
			)}
			{...props}
		>
			<svg
				width={size}
				height={size}
				viewBox="0 0 24 24"
				fill="none"
				aria-hidden="true"
			>
				<rect
					x="2.5"
					y="2.5"
					width="19"
					height="19"
					rx="3.5"
					stroke="currentColor"
					strokeWidth="1.6"
				/>
				<rect
					x="7"
					y="7"
					width="10"
					height="10"
					rx="1.5"
					className="stroke-primary"
					strokeWidth="1.6"
				/>
				<rect x="10.5" y="10.5" width="3" height="3" className="fill-primary" />
				<path
					d="M2.5 12h2M19.5 12h2M12 2.5v2M12 19.5v2"
					stroke="currentColor"
					strokeWidth="1.6"
					strokeLinecap="round"
				/>
			</svg>
			{withWord && (
				<span
					className="font-semibold tracking-tight"
					style={{ fontSize: size * 0.78 }}
				>
					Local<span className="text-primary">Vault</span>
				</span>
			)}
		</span>
	);
}

export { LVLogo };
