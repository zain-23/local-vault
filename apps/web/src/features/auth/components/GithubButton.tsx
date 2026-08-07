import type * as React from "react";

import { GithubIcon } from "#/components/shared";
import { Button } from "#/components/ui";
import { authService } from "../api";
import { setPostLoginRedirect } from "../utils/postLoginRedirect.ts";

// The only sign-in method: first GitHub login creates the account, a
// returning one logs in — there's no separate signup flow to keep in sync.
function GithubButton({
	children = "Continue with GitHub",
	redirectTo,
	...props
}: React.ComponentProps<typeof Button> & { redirectTo?: string }) {
	const onClick = () => {
		if (redirectTo) setPostLoginRedirect(redirectTo);
		window.location.href = authService.oauthUrl();
	};
	return (
		<Button
			type="button"
			variant="secondary"
			className="w-full"
			size={"lg"}
			onClick={onClick}
			{...props}
		>
			<GithubIcon size={22}/>
			{children}
		</Button>
	);
}

export { GithubButton };
