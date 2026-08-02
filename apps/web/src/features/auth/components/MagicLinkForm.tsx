import { zodResolver } from "@hookform/resolvers/zod";
import { Link, useNavigate } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldGroup,
	FieldLabel,
	Input,
} from "#/components/ui/index.ts";
import { useSendMagicLink } from "../hooks/useSendMagicLink.ts";
import { type MagicLinkValues, magicLinkSchema } from "../schemas/index.ts";
import { AuthFooter } from "./AuthFooter.tsx";
import { AuthHeading } from "./AuthHeading.tsx";

function MagicLinkForm() {
	const navigate = useNavigate();
	const { mutate, isPending } = useSendMagicLink();
	const {
		register,
		handleSubmit,
		formState: { errors },
	} = useForm<MagicLinkValues>({
		resolver: zodResolver(magicLinkSchema),
		defaultValues: { email: "" },
	});

	const onSubmit = handleSubmit(({ email }) => {
		mutate(email, {
			onSuccess: () => {
				navigate({
					to: "/auth/check-email",
					search: { variant: "magic", email },
				});
			},
		});
	});

	return (
		<>
			<AuthHeading
				title="Sign in with a magic link"
				subtitle="We'll email you a link that signs you in instantly — no password to remember."
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
						Email me a link
						<ArrowRight />
					</Button>
				</FieldGroup>
			</form>

			<AuthFooter>
				Prefer a password?{" "}
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

export { MagicLinkForm };
