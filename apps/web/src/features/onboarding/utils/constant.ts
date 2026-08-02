import { Play, Users } from "lucide-react";

export const PLATFORMS = [
	{ key: "macos", label: "macOS", cmd: "brew install localvault/tap/lv" },
	{
		key: "linux",
		label: "Linux",
		cmd: "curl -fsSL https://localvault.app/install.sh | sh",
	},
	{
		key: "source",
		label: "From source",
		cmd: "go install github.com/localvault/cli@latest",
	},
] as const;

// Step 4 — first commands to try once onboarding is done (UI-only prompts).
export const NEXT_STEPS = [
	{
		icon: Play,
		label: "Inject secrets into a process",
		command: "lv run -- npm start",
	},
	{
		icon: Users,
		label: "Invite a teammate",
		command: "lv team invite dev@kodexo-labs.com",
	},
] as const;
