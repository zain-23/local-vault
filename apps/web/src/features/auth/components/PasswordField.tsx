import { Eye, EyeOff } from "lucide-react";
import * as React from "react";

import { Input } from "#/components/ui/index.ts";
import { cn } from "#/lib/utils.ts";

// Password input with an inline show/hide toggle. All extra props (including
// react-hook-form's `register(...)` output — value, onChange, onBlur, name, ref)
// are forwarded straight to the native input, so it plugs into RHF unchanged.
function PasswordField({
	className,
	...props
}: React.ComponentProps<typeof Input>) {
	const [show, setShow] = React.useState(false);

	return (
		<div className="relative">
			<Input
				type={show ? "text" : "password"}
				className={cn("pr-9", className)}
				{...props}
			/>
			<button
				type="button"
				// preventDefault keeps focus in the input when toggling visibility.
				onMouseDown={(e) => e.preventDefault()}
				onClick={() => setShow((s) => !s)}
				className="absolute inset-y-0 right-0 flex w-9 items-center justify-center text-muted-foreground transition-colors hover:text-foreground"
				title={show ? "Hide password" : "Show password"}
				aria-label={show ? "Hide password" : "Show password"}
			>
				{show ? <EyeOff className="size-4" /> : <Eye className="size-4" />}
			</button>
		</div>
	);
}

export { PasswordField };
