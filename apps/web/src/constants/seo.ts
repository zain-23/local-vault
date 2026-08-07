import type { PageMeta } from "#/types";
// Direct path, not "#/constants": this file is re-exported by that barrel, so
// going through it would be a cycle and leave these consts in the TDZ.
import { DEFAULT_DESCRIPTION, DEFAULT_TITLE, SITE_NAME } from "./site.ts";

/** Keeps the app, auth flows and device-link screens out of search results. */
export const PRIVATE_ROBOTS = "noindex, nofollow";

/** `"Sign in" → "Sign in — LocalVault"`. Local formatter, not public API. */
const titled = (page: string) => `${page} — ${SITE_NAME}`;

/** Site-wide fallback, applied on the root route and overridden per page. */
export const DEFAULT_PAGE_META: PageMeta = {
	title: DEFAULT_TITLE,
	description: DEFAULT_DESCRIPTION,
};

/**
 * One entry per route, keyed by the path passed to `createFileRoute`. Add the
 * route here, then `head: () => seo(PAGE_META["/your/route"])` in the route file
 * — indexing is type-checked, so a typo or a missing entry fails the build.
 */
export const PAGE_META = {
	"/": {
		title: DEFAULT_TITLE,
		description: DEFAULT_DESCRIPTION,
		path: "/",
	},

	/* ── auth ── */
	"/auth/login": {
		title: titled("Sign in"),
		description: "Sign in to your LocalVault workspace with GitHub.",
		robots: PRIVATE_ROBOTS,
	},

	/* ── onboarding & joining ── */
	"/onboarding": {
		title: titled("Set up your workspace"),
		robots: PRIVATE_ROBOTS,
	},
	"/workspaces/join": {
		title: titled("Join a workspace"),
		robots: PRIVATE_ROBOTS,
	},

	/* ── CLI device pairing ── */
	"/device/": {
		title: titled("Link a device"),
		robots: PRIVATE_ROBOTS,
	},
	"/device/confirmation": {
		title: titled("Device linked"),
		robots: PRIVATE_ROBOTS,
	},

	/* ── the app itself ── */
	"/_app/dashboard": { title: titled("Dashboard"), robots: PRIVATE_ROBOTS },
	"/_app/vaults/": { title: titled("Vaults"), robots: PRIVATE_ROBOTS },
	"/_app/vaults/$vaultId": { title: titled("Vault"), robots: PRIVATE_ROBOTS },
	"/_app/members": { title: titled("Members"), robots: PRIVATE_ROBOTS },
	"/_app/settings": { title: titled("Settings"), robots: PRIVATE_ROBOTS },
	"/_app/audit": { title: titled("Audit log"), robots: PRIVATE_ROBOTS },
} satisfies Record<string, PageMeta>;
