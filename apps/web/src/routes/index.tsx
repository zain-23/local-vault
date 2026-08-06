import { createFileRoute } from "@tanstack/react-router";

import { LandingPage } from "#/features/marketing/components/LandingPage.tsx";

// Public marketing landing. Signed-in users reach the app from the nav; the
// guard on `/_app` still owns everything behind it.
export const Route = createFileRoute("/")({
	head: () => ({
		meta: [
			{ title: "LocalVault — Stop sharing secrets over Slack" },
			{
				name: "description",
				content:
					"LocalVault replaces .env files with an encrypted vault that syncs peer-to-peer between your team's machines. No cloud storage, works offline, MIT licensed.",
			},
		],
	}),
	component: LandingPage,
});
