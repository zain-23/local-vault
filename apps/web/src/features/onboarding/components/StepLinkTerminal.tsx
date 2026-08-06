import { ArrowRight } from "lucide-react";

import { Button } from "#/components/ui/index.ts";
import { CodeBlock } from "./CodeBlock.tsx";
import { StepHeading } from "./OnboardingLayout.tsx";

// Step 3 — link a terminal. The "waiting" card is purely visual; the user
// advances with an explicit Continue button (no fake auto-advance).
function StepLinkTerminal({ onContinue }: { onContinue: () => void }) {
	return (
		<>
			<StepHeading
				title="Link your terminal"
				subtitle="Run this command on the machine you want to authorize."
			/>

			<CodeBlock command="lv login --workspace kodexo-labs" />

			<Button className="mt-6 w-full" onClick={onContinue}>
				Continue
				<ArrowRight className="size-3.5" />
			</Button>
		</>
	);
}

export { StepLinkTerminal };
