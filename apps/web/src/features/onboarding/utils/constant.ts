import { Play, Users } from "lucide-react";

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
