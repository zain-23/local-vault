import { zodResolver } from "@hookform/resolvers/zod";
import { Link } from "@tanstack/react-router";
import { useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "#/components/ui/index.ts";
import { useResetPassword } from "../hooks/useResetPassword.ts";
import {
	type ResetPasswordValues,
	resetPasswordSchema,
} from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { PasswordField } from "./PasswordField.tsx";
import { PasswordStrengthMeter } from "./PasswordStrengthMeter.tsx";

function InvalidResetLink() {
	return (
		<>
			<AuthHeading
				title="Reset link invalid"
				subtitle="This password reset link is missing, expired, or already used. Request a new one to continue."
			/>

			<Button asChild className="w-full">
				<Link to="/auth/forgot-password">Request a new link</Link>
			</Button>

			<AuthFooter>
				Back to{" "}
				<Link
					to="/auth/login"
					className="font-medium text-foreground hover:underline"
				>
					Log in
				</Link>
			</AuthFooter>
		</>
	);
}

function ResetPasswordForm({ token }: { token?: string }) {
	const { mutate, isPending, isError } = useResetPassword();
	const {
		register,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<ResetPasswordValues>({
		resolver: zodResolver(resetPasswordSchema),
		defaultValues: { password: "", confirmPassword: "" },
	});

	const password = watch("password");

	// Missing from the URL, or rejected by the API (expired / already used / unknown).
	if (!token || isError) {
		return <InvalidResetLink />;
	}

	const onSubmit = handleSubmit(({ password: newPassword }) => {
		mutate({ token, new_password: newPassword });
	});

	return (
		<>
			<AuthHeading
				title="Set a new password"
				subtitle="Choose a strong password you haven't used on LocalVault before."
			/>

			<form onSubmit={onSubmit} noValidate>
				<FieldGroup className="gap-4">
					<Field data-invalid={!!errors.password}>
						<FieldLabel htmlFor="password">New password</FieldLabel>
						<PasswordField
							id="password"
							autoComplete="new-password"
							placeholder="At least 12 characters"
							aria-invalid={!!errors.password}
							disabled={isPending}
							{...register("password")}
						/>
						<PasswordStrengthMeter value={password} />
						<FieldError errors={[errors.password]} />
					</Field>

					<Field data-invalid={!!errors.confirmPassword}>
						<FieldLabel htmlFor="confirmPassword">Confirm password</FieldLabel>
						<PasswordField
							id="confirmPassword"
							autoComplete="new-password"
							placeholder="Re-enter your password"
							aria-invalid={!!errors.confirmPassword}
							disabled={isPending}
							{...register("confirmPassword")}
						/>
						<FieldError errors={[errors.confirmPassword]} />
					</Field>

					<Button
						type="submit"
						className="w-full"
						disabled={isPending}
						isLoading={isPending}
					>
						Update password
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Back to{" "}
				<Link
					to="/auth/login"
					className="font-medium text-foreground hover:underline"
				>
					Log in
				</Link>
			</AuthFooter>
		</>
	);
}

export { ResetPasswordForm };
