import { Link } from "@tanstack/react-router";
import { ArrowRight, Check } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { AuthHeading } from "./AuthHeading.tsx";

// Terminal confirmation after a successful password change.
function ResetSuccessPanel() {
	return (
		<>
			<AuthHeading
				icon={<Check className="size-5.5" strokeWidth={2.5} />}
				iconTone="success"
				title="Password updated"
				subtitle="Your password has been changed. Sign in with your new password to continue."
			/>

			<Button asChild size="lg" className="w-full">
				<Link to="/auth/login">
					Continue to login
					<ArrowRight />
				</Link>
			</Button>
		</>
	);
}

export { ResetSuccessPanel };
