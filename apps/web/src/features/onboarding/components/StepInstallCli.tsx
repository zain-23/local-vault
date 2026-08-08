import { Button } from "#/components/ui/index.ts";
import { InstallCommand } from "#/features/marketing/components/InstallCommand.tsx";
import { StepHeading } from "./OnboardingLayout.tsx";

function StepInstallCli({ onContinue }: { onContinue: () => void }) {
	return (
		<>
			<StepHeading
				title="Install the CLI"
				subtitle="LocalVault is CLI-first. Install once per machine."
			/>

			{/* Same install command as the marketing page, so the two never drift apart. */}
			<div className="flex flex-col items-center">
				<InstallCommand />
			</div>

			<div className="mt-6 flex gap-2">
				<Button variant="outline" className="flex-1" onClick={onContinue}>
					Skip for now
				</Button>
				<Button className="flex-1" onClick={onContinue}>
					I've installed it
				</Button>
			</div>
		</>
	);
}

export { StepInstallCli };
