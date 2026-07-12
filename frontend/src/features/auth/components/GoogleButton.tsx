import type * as React from "react";
import { Button } from "#/components/ui";

// Social sign-in button, shared by login + signup. Full width, no real auth —
// this is UI only, so the click is a no-op until integration is wired up.
// (lucide-react ships no Google brand mark, so this stays text-only.)
function GoogleButton({
	children = "Continue with Google",
	...props
}: React.ComponentProps<typeof Button>) {
	return (
		<Button
			type="button"
			variant="secondary"
			size="lg"
			className="w-full"
			{...props}
		>
			{children}
		</Button>
	);
}

export { GoogleButton };
