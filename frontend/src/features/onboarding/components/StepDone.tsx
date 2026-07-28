import { Link } from "@tanstack/react-router";
import { ArrowRight, Check } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { StepHeading } from "./OnboardingLayout.tsx";

// Summary of what onboarding set up (mock values, UI only).
const SUMMARY = [
	{ label: "Workspace created", value: "kodexo-labs" },
	{ label: "CLI installed", value: "v1.0.2" },
	{ label: "Device linked", value: "ahmed-mbp · b427d8d8…" },
];

// Step 4 — completion. CTA links to the app root.
function StepDone() {
	return (
		<>
			<StepHeading
				title="You're all set"
				subtitle="Your workspace is ready. Here's what to try next."
			/>

			<div className="flex flex-col gap-2">
				{SUMMARY.map((item) => (
					<div
						key={item.label}
						className="flex items-center gap-2.5 rounded-md border border-border px-3 py-2.5"
					>
						<span className="flex size-[18px] items-center justify-center rounded-full bg-success text-white">
							<Check className="size-3" strokeWidth={2.5} />
						</span>
						<span className="flex-1 text-[13px] font-medium">{item.label}</span>
						<span className="font-mono text-[11.5px] text-muted-foreground">
							{item.value}
						</span>
					</div>
				))}
			</div>

			<Button size="lg" className="mt-6 w-full" asChild>
				<Link to="/">
					Go to dashboard
					<ArrowRight className="size-3.5" />
				</Link>
			</Button>
		</>
	);
}

export { StepDone };
