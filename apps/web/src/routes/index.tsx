import { createFileRoute } from "@tanstack/react-router";

import { PAGE_META } from "#/constants";
import { LandingPage } from "#/features/marketing/components/LandingPage.tsx";
import { seo } from "#/utils/seo.ts";

// Public marketing landing. Signed-in users reach the app from the nav; the
// guard on `/_app` still owns everything behind it.
export const Route = createFileRoute("/")({
	head: () => seo(PAGE_META["/"]),
	component: LandingPage,
});
