import { cn } from "#/lib/utils.ts";

// Cheap, dependency-free strength heuristic (0–4) based on length + character
// variety. Enough to give the four-bar meter meaningful feedback; real strength
// estimation (zxcvbn etc.) can slot in later without changing this component.
function scorePassword(value: string): number {
	if (!value) return 0;
	let score = 0;
	if (value.length >= 8) score++;
	if (value.length >= 12) score++;
	if (/[0-9]/.test(value) && /[a-z]/.test(value) && /[A-Z]/.test(value))
		score++;
	if (/[^A-Za-z0-9]/.test(value)) score++;
	return Math.min(score, 4);
}

const LABELS = ["", "Weak", "Fair", "Good", "Strong"] as const;
// Bar colors per filled segment, indexed by the meter's score.
const TONES = [
	"bg-destructive",
	"bg-warning",
	"bg-warning",
	"bg-success",
] as const;

// Four-segment meter. Renders nothing until the user starts typing.
function PasswordStrengthMeter({ value }: { value: string }) {
	const score = scorePassword(value);
	if (!value) return null;

	return (
		<div className="flex items-center gap-1.5 text-[11.5px] text-muted-foreground">
			{[0, 1, 2, 3].map((i) => (
				<span
					key={i}
					className={cn(
						"h-[3px] flex-1 rounded-full",
						i < score ? TONES[score - 1] : "bg-border-strong",
					)}
				/>
			))}
			<span className="ml-1 tabular-nums">{LABELS[score]}</span>
		</div>
	);
}

export { PasswordStrengthMeter, scorePassword };
