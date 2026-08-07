import { useSearch } from "@tanstack/react-router";

import { AuthHeading } from "./AuthHeading.tsx";
import { GithubButton } from "./GithubButton.tsx";
import { authErrorMessage } from "../utils/authErrorMessage.ts";

function LoginForm() {
	// Read loosely: this component only renders on /auth/login, but strict:false
	// keeps it decoupled from that route's typed search.
	const search = useSearch({ strict: false }) as {
		error?: string;
		redirect?: string;
	};

	return (
		<>
			<AuthHeading
				title="Welcome back"
				subtitle="Sign in to your LocalVault workspaces with GitHub."
			/>

			<GithubButton redirectTo={search.redirect} />

			{search.error && (
				<p className="mt-4 text-center text-sm text-destructive">
					{authErrorMessage(search.error)}
				</p>
			)}
		</>
	);
}

export { LoginForm };
