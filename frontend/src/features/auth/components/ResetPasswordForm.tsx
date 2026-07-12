import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
} from "#/components/ui/index.ts";
import {
	type ResetPasswordValues,
	resetPasswordSchema,
} from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { PasswordField } from "./PasswordField.tsx";
import { PasswordStrengthMeter } from "./PasswordStrengthMeter.tsx";

function ResetPasswordForm() {
	const navigate = useNavigate();
	const {
		register,
		watch,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<ResetPasswordValues>({
		resolver: zodResolver(resetPasswordSchema),
		defaultValues: { password: "", confirmPassword: "" },
	});

	const password = watch("password");
	const onSubmit = handleSubmit(() => navigate({ to: "/auth/reset-success" }));

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
							{...register("confirmPassword")}
						/>
						<FieldError errors={[errors.confirmPassword]} />
					</Field>

					<Button
						type="submit"
						size="lg"
						className="w-full"
						disabled={isSubmitting}
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
