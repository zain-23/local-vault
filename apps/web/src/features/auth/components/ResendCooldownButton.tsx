import { Button } from "#/components/ui/index.ts";

const RING_RADIUS = 7;
const RING_CIRCUMFERENCE = 2 * Math.PI * RING_RADIUS;

function formatCountdown(seconds: number) {
	const m = Math.floor(seconds / 60);
	const s = seconds % 60;
	return `${m}:${String(s).padStart(2, "0")}`;
}

function CountdownRing({ progress }: { progress: number }) {
	return (
		<svg
			aria-hidden="true"
			viewBox="0 0 20 20"
			className="size-4 -rotate-90"
			focusable="false"
		>
			<title>Time remaining</title>
			<circle
				cx="10"
				cy="10"
				r={RING_RADIUS}
				fill="none"
				stroke="currentColor"
				strokeWidth="2"
				className="opacity-25"
			/>
			<circle
				cx="10"
				cy="10"
				r={RING_RADIUS}
				fill="none"
				stroke="currentColor"
				strokeWidth="2"
				strokeLinecap="round"
				strokeDasharray={RING_CIRCUMFERENCE}
				strokeDashoffset={RING_CIRCUMFERENCE * (1 - progress)}
				className="transition-[stroke-dashoffset] duration-75 ease-linear"
			/>
		</svg>
	);
}

function ResendCooldownButton({
	isCoolingDown,
	isPending,
	secondsLeft,
	progress,
	onResend,
}: {
	isCoolingDown: boolean;
	isPending: boolean;
	secondsLeft: number;
	progress: number;
	onResend: () => void;
}) {
	const disabled = isCoolingDown || isPending;

	return (
		<Button
			type="button"
			className="w-full"
			disabled={disabled}
			isLoading={isPending && !isCoolingDown}
			onClick={onResend}
			aria-live="polite"
		>
			{isCoolingDown ? (
				<>
					<CountdownRing progress={progress} />
					Resend in {formatCountdown(secondsLeft)}
				</>
			) : (
				<>Resend email</>
			)}
		</Button>
	);
}

export { ResendCooldownButton };
