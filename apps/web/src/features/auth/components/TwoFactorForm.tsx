import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import { Controller, useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
} from "#/components/ui/index.ts";
import { useLogin2FA } from "../hooks/useLogin2FA.ts";
import { type TwoFactorValues, twoFactorSchema } from "../schemas/index.ts";
import { clearTempToken, getTempToken } from "../utils/tempToken.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { OtpField } from "./OtpField.tsx";

function TwoFactorForm() {
	const navigate = useNavigate();
	const { mutate, isPending } = useLogin2FA();
	const {
		control,
		handleSubmit,
		formState: { errors },
	} = useForm<TwoFactorValues>({
		resolver: zodResolver(twoFactorSchema),
		defaultValues: { code: "" },
	});

	// No temp token means the user landed here cold — send them back to login.
	useEffect(() => {
		if (!getTempToken()) {
			navigate({ to: "/auth/login", replace: true });
		}
	}, [navigate]);

	const onSubmit = handleSubmit(({ code }) => {
		mutate(code);
	});

	return (
		<>
			<AuthHeading
				title="Two-factor authentication"
				subtitle="Enter the 6-digit code from your authenticator app."
			/>

			<form onSubmit={onSubmit} noValidate>
				<FieldGroup className="gap-4">
					<Field data-invalid={!!errors.code}>
						<Controller
							control={control}
							name="code"
							render={({ field }) => (
								<OtpField
									value={field.value}
									onChange={field.onChange}
									// Auto-submit once all six digits are entered.
									onComplete={(value) => mutate(value)}
									aria-invalid={!!errors.code}
									disabled={isPending}
								/>
							)}
						/>
						<FieldError errors={[errors.code]} />
					</Field>

					<Button
						type="submit"
						className="w-full"
						disabled={isPending}
						isLoading={isPending}
					>
						Verify
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Wrong account?{" "}
				<Link
					to="/auth/login"
					onClick={() => clearTempToken()}
					className="font-medium text-foreground hover:underline"
				>
					Back to login
				</Link>
			</AuthFooter>
		</>
	);
}

export { TwoFactorForm };
