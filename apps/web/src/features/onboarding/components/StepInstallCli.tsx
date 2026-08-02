import {
	Button,
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from "#/components/ui/index.ts";
import { PLATFORMS } from "../utils/constant.ts";
import { CodeBlock } from "./CodeBlock.tsx";
import { StepHeading } from "./OnboardingLayout.tsx";

function StepInstallCli({ onContinue }: { onContinue: () => void }) {
	return (
		<>
			<StepHeading
				title="Install the CLI"
				subtitle="LocalVault is CLI-first. Install once per machine."
			/>

			{/* platform tabs; Radix tracks the active tab and swaps the command */}
			<Tabs defaultValue="linux">
				<TabsList className="w-full">
					{PLATFORMS.map((p) => (
						<TabsTrigger key={p.key} value={p.key} className="text-xs">
							{p.label}
						</TabsTrigger>
					))}
				</TabsList>
				{PLATFORMS.map((p) => (
					<TabsContent key={p.key} value={p.key}>
						<CodeBlock command={p.cmd} />
					</TabsContent>
				))}
			</Tabs>

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
