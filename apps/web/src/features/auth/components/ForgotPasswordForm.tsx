import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
	Input,
} from "#/components/ui/index.ts";
import { useForgotPassword } from "../hooks/useForgotPassword.ts";
import {
	type ForgotPasswordValues,
	forgotPasswordSchema,
} from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";

function ForgotPasswordForm() {
	const navigate = useNavigate();
	const { mutate, isPending } = useForgotPassword();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<ForgotPasswordValues>({
		resolver: zodResolver(forgotPasswordSchema),
		defaultValues: { email: "" },
	});

	const onSubmit = handleSubmit(({ email }) => {
		mutate(email, {
			onSuccess: () => {
				navigate({
					to: "/auth/check-email",
					search: { variant: "reset", email },
				});
			},
		});
	});

	return (
		<>
			<AuthHeading
				title="Reset your password"
				subtitle="Enter your account email and we'll send you a link to set a new password."
			/>

			<form onSubmit={onSubmit} noValidate>
				<FieldGroup className="gap-4">
					<Field data-invalid={!!errors.email}>
						<FieldLabel htmlFor="email">Email</FieldLabel>
						<Input
							id="email"
							type="email"
							autoComplete="email"
							placeholder="ada@company.com"
							aria-invalid={!!errors.email}
							disabled={isPending}
							{...register("email")}
						/>
						<FieldError errors={[errors.email]} />
					</Field>

					<Button
						type="submit"
						className="w-full"
						disabled={isPending}
						isLoading={isPending}
					>
						Send reset link
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Remembered it?{" "}
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

export { ForgotPasswordForm };
