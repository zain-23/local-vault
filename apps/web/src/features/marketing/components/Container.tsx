import type * as React from "react";

import { cn } from "#/lib/utils.ts";

// The single content measure the whole landing page is built on.
function Container({ className, ...props }: React.ComponentProps<"div">) {
	return (
		<div
			className={cn(
				"mx-auto w-full max-w-[1180px] px-[18px] sm:px-6",
				className,
			)}
			{...props}
		/>
	);
}

export { Container };
