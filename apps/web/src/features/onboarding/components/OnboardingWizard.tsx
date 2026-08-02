import { useState } from "react";

import { OnboardingLayout } from "./OnboardingLayout.tsx";
import { StepDone } from "./StepDone.tsx";
import { StepInstallCli } from "./StepInstallCli.tsx";
import { StepLinkTerminal } from "./StepLinkTerminal.tsx";
import { StepWorkspace } from "./StepWorkspace.tsx";

const TOTAL_STEPS = 4;

// Container for the post-signup wizard. Holds the current step and advances
// forward only — no back navigation, no integration. Each step gets onContinue.
function OnboardingWizard() {
	const [step, setStep] = useState(1);
	const goNext = () => setStep((s) => Math.min(s + 1, TOTAL_STEPS));

	return (
		<OnboardingLayout step={step}>
			{step === 1 && <StepWorkspace onContinue={goNext} />}
			{step === 2 && <StepInstallCli onContinue={goNext} />}
			{step === 3 && <StepLinkTerminal onContinue={goNext} />}
			{step === 4 && <StepDone />}
		</OnboardingLayout>
	);
}

export { OnboardingWizard };
