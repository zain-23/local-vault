import type * as React from "react";

import { GoogleIcon } from "#/components/shared/GoogleIcon.tsx";
import { Button } from "#/components/ui";
import { authService } from "../api";

// Social sign-in button, shared by login + signup. Full width, no real auth —
// this is UI only, so the click is a no-op until integration is wired up.
function GoogleButton({
	children = "Continue with Google",
	...props
}: React.ComponentProps<typeof Button>) {
	const onClick = () => {
		window.location.href = authService.oauthUrl();
	};
	return (
		<Button
			type="button"
			variant="secondary"
			className="w-full"
			onClick={onClick}
			{...props}
		>
			<GoogleIcon />
			{children}
		</Button>
	);
}

export { GoogleButton };
