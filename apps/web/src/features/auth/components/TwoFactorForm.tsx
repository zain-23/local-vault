import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
} from "#/components/ui/index.ts";
import { type TwoFactorValues, twoFactorSchema } from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { OtpField } from "./OtpField.tsx";

function TwoFactorForm() {
	const navigate = useNavigate();
	const {
		control,
		handleSubmit,
		formState: { errors, isSubmitting },
	} = useForm<TwoFactorValues>({
		resolver: zodResolver(twoFactorSchema),
		defaultValues: { code: "" },
	});

	const onSubmit = handleSubmit(() => navigate({ to: "/" }));

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
									onComplete={() => onSubmit()}
									aria-invalid={!!errors.code}
								/>
							)}
						/>
						<FieldError errors={[errors.code]} />
					</Field>

					<Button type="submit" className="w-full" disabled={isSubmitting}>
						Verify
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Lost your device?{" "}
				<Link
					to="/auth/login"
					className="font-medium text-foreground hover:underline"
				>
					Use a backup code
				</Link>
			</AuthFooter>
		</>
	);
}

export { TwoFactorForm };
