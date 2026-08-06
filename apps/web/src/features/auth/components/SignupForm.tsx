import { zodResolver } from "@hookform/resolvers/zod";
import { Link } from "@tanstack/react-router";
import { Controller, useForm } from "react-hook-form";

import {
	Button,
	Checkbox,
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
	FieldSeparator,
	Input,
} from "#/components/ui/index.ts";
import { useSignup } from "../hooks/useSignup.ts";
import { type SignupValues, signupSchema } from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";
import { GoogleButton } from "./GoogleButton.tsx";
import { PasswordField } from "./PasswordField.tsx";
import { PasswordStrengthMeter } from "./PasswordStrengthMeter.tsx";

function SignupForm() {
	const {
		register,
		control,
		watch,
		handleSubmit,
		formState: { errors },
	} = useForm<SignupValues>({
		resolver: zodResolver(signupSchema),
		defaultValues: { name: "", email: "", password: "", terms: false },
	});

	const { isPending, mutate } = useSignup();

	const password = watch("password");

	const onSubmit = handleSubmit((values) => {
		mutate(values);
	});

	return (
		<>
			<AuthHeading
				title="Create your account"
				subtitle="Get started with LocalVault. No credit card."
			/>

			<GoogleButton disabled={isPending} />

			<FieldSeparator className="my-5">OR</FieldSeparator>

			<form onSubmit={onSubmit} noValidate>
				<FieldGroup className="gap-4">
					<Field data-invalid={!!errors.name}>
						<FieldLabel htmlFor="name">Name</FieldLabel>
						<Input
							id="name"
							autoComplete="name"
							placeholder="Ada Lovelace"
							aria-invalid={!!errors.name}
							{...register("name")}
						/>
						<FieldError errors={[errors.name]} />
					</Field>

					<Field data-invalid={!!errors.email}>
						<FieldLabel htmlFor="email">Work email</FieldLabel>
						<Input
							id="email"
							type="email"
							autoComplete="email"
							placeholder="ada@company.com"
							aria-invalid={!!errors.email}
							{...register("email")}
						/>
						<FieldError errors={[errors.email]} />
					</Field>

					<Field data-invalid={!!errors.password}>
						<FieldLabel htmlFor="password">Password</FieldLabel>
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

					<Field data-invalid={!!errors.terms}>
						<Controller
							control={control}
							name="terms"
							render={({ field }) => (
								<label
									htmlFor="terms"
									className="flex items-start gap-2 text-[12.5px] text-muted-foreground"
								>
									<Checkbox
										id="terms"
										checked={field.value}
										onCheckedChange={(v) => field.onChange(v === true)}
										onBlur={field.onBlur}
										aria-invalid={!!errors.terms}
										className="mt-0.5"
									/>
									<span>
										I agree to the{" "}
										<Link
											to="/auth/signup"
											className="text-foreground hover:underline"
										>
											Terms
										</Link>{" "}
										and{" "}
										<Link
											to="/auth/signup"
											className="text-foreground hover:underline"
										>
											Privacy Policy
										</Link>
									</span>
								</label>
							)}
						/>
						<FieldError errors={[errors.terms]} />
					</Field>

					<Button type="submit" className="w-full" isLoading={isPending}>
						Create account
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Already have an account?{" "}
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

export { SignupForm };
