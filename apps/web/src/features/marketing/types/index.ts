import type { LucideIcon } from "lucide-react";

/** Anchor or outbound link in the nav / footer. */
export type MarketingLink = {
	label: string;
	href: string;
	/** Leaves the site — opens in a new tab. */
	external?: boolean;
};

export type FooterColumn = {
	id: string;
	title: string;
	links: MarketingLink[];
};

/** Short "no cloud storage / works offline" claim under the hero CTAs. */
export type TrustPoint = {
	id: string;
	label: string;
};

/* ── product tour ───────────────────────────────────────────────────── */

export type ProductTabId = "cli" | "web";

export type ProductTab = {
	id: ProductTabId;
	label: string;
	/** Text in the window chrome while this tab is active. */
	frameTitle: string;
	/** Badge on the right of the window chrome. */
	frameBadge: string;
};

/**
 * Color role for a run of terminal text. Keeps the transcript plain data so it
 * can live in `utils/constants.ts` instead of leaking JSX into it.
 */
export type TerminalTone = "prompt" | "success" | "warning" | "muted";

export type TerminalSegment = {
	text: string;
	tone?: TerminalTone;
};

export type TerminalLine = {
	id: string;
	segments: TerminalSegment[];
	/** Renders the blinking block cursor at the end of the line. */
	caret?: boolean;
};

export type ConsoleNavItem = {
	id: string;
	label: string;
	icon: LucideIcon;
	active?: boolean;
};

export type ConsoleStat = {
	id: string;
	label: string;
	value: number;
	/** Counts up from zero on mount; otherwise the value renders as-is. */
	countUp?: boolean;
	tone?: "default" | "success";
};

export type VaultSyncStatus = "synced" | "pending";

export type VaultRow = {
	id: string;
	name: string;
	secrets: number;
	members: number;
	/** Humanised relative time, e.g. "2m ago". */
	lastSync: string;
	status: VaultSyncStatus;
	statusLabel: string;
};

/* ── content sections ───────────────────────────────────────────────── */

export type OnboardingStep = {
	id: string;
	title: string;
	description: string;
	command: string;
};

export type Feature = {
	id: string;
	icon: LucideIcon;
	title: string;
	description: string;
	/** Mono footer line, e.g. `$ lv pair`. */
	command?: string;
	/** Emphasised token inside `command`. */
	commandAccent?: string;
};

/** One row of "what the relay can see". */
export type RelayLedgerEntry = {
	id: string;
	key: string;
	value: string;
	/** True when the relay genuinely observes this, false when it never sees it. */
	exposed: boolean;
};

export type CryptoPrimitive = {
	id: string;
	name: string;
	role: string;
	description: string;
};
