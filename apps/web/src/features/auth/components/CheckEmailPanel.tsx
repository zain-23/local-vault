import { Link } from "@tanstack/react-router";

import { useForgotPassword } from "../hooks/useForgotPassword.ts";
import { useResendCooldown } from "../hooks/useResendCooldown.ts";
import { useSendMagicLink } from "../hooks/useSendMagicLink.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { ResendCooldownButton } from "./ResendCooldownButton.tsx";

// Shared "check your inbox" confirmation for both the password-reset and
// magic-link flows — the copy switches on `variant`.
function CheckEmailPanel({
	variant,
	email,
}: {
	variant: "reset" | "magic";
	email?: string;
}) {
	const forgotPassword = useForgotPassword();
	const sendMagicLink = useSendMagicLink();
	const { secondsLeft, progress, isCoolingDown, restart } = useResendCooldown();

	const resend = variant === "magic" ? sendMagicLink : forgotPassword;
	const canResend = Boolean(email);

	const subtitle =
		variant === "magic" ? (
			<>
				We sent a magic sign-in link to{" "}
				<strong className="text-foreground">{email ?? "your email"}</strong>.
				Open it on this device to continue.
			</>
		) : (
			<>
				We sent a password reset link to{" "}
				<strong className="text-foreground">{email ?? "your email"}</strong>.
				The link expires in 15 minutes.
			</>
		);

	const handleResend = () => {
		if (!email || resend.isPending || isCoolingDown) return;
		resend.mutate(email, {
			onSuccess: () => restart(),
		});
	};

	return (
		<>
			<AuthHeading title="Check your inbox" subtitle={subtitle} />

			{canResend ? (
				<ResendCooldownButton
					isCoolingDown={isCoolingDown}
					isPending={resend.isPending}
					secondsLeft={secondsLeft}
					progress={progress}
					onResend={handleResend}
				/>
			) : null}

			<AuthFooter>
				Wrong address?{" "}
				<Link
					to={
						variant === "magic" ? "/auth/magic-link" : "/auth/forgot-password"
					}
					className="font-medium text-foreground hover:underline"
				>
					Use a different one
				</Link>
			</AuthFooter>
		</>
	);
}

export { CheckEmailPanel };
