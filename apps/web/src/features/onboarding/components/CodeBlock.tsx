import { Check, Copy } from "lucide-react";
import { useState } from "react";

import { cn } from "#/lib/utils.ts";

function CodeBlock({ command }: { command: string }) {
	const [copied, setCopied] = useState(false);

	const copy = () => {
		// navigator.clipboard is async and may be unavailable in insecure contexts.
		navigator.clipboard?.writeText(command).then(() => {
			setCopied(true);
			setTimeout(() => setCopied(false), 1500);
		});
	};

	return (
		<div className="flex items-center gap-3 rounded-lg border border-border bg-muted px-3.5 py-3 font-mono text-[13px]">
			<span className="text-muted-foreground select-none">$</span>
			<code className="flex-1 overflow-x-auto whitespace-nowrap text-foreground">
				{command}
			</code>
			<button
				type="button"
				onClick={copy}
				aria-label="Copy command"
				className={cn(
					"shrink-0 rounded-md p-1 transition-colors hover:text-foreground",
					copied ? "text-success" : "text-muted-foreground",
				)}
			>
				{copied ? (
					<Check className="size-3.5" />
				) : (
					<Copy className="size-3.5" />
				)}
			</button>
		</div>
	);
}

export { CodeBlock };
