import {
	InputOTP,
	InputOTPGroup,
	InputOTPSlot,
} from "#/components/ui/index.ts";

// Six-digit OTP entry for the two-factor screen. Controlled (value/onChange) so
// it drops into a react-hook-form <Controller/>. Slots are enlarged and mono to
// match the LocalVault auth styling.
function OtpField({
	value,
	onChange,
	onComplete,
	disabled,
	"aria-invalid": ariaInvalid,
}: {
	value: string;
	onChange: (value: string) => void;
	onComplete?: (value: string) => void;
	disabled?: boolean;
	"aria-invalid"?: boolean;
}) {
	return (
		<InputOTP
			maxLength={6}
			value={value}
			onChange={onChange}
			onComplete={onComplete}
			disabled={disabled}
			containerClassName="justify-between gap-2"
		>
			<InputOTPGroup className="w-full gap-2">
				{[0, 1, 2, 3, 4, 5].map((i) => (
					<InputOTPSlot
						key={i}
						index={i}
						aria-invalid={ariaInvalid}
						className="h-12 flex-1 rounded-md! border-l! font-mono text-lg font-semibold"
					/>
				))}
			</InputOTPGroup>
		</InputOTP>
	);
}

export { OtpField };
