import {
	debounce,
	parseAsInteger,
	parseAsString,
	parseAsStringLiteral,
} from "nuqs";

import { MEMBER_ROLES } from "#/features/members/utils";

export { MEMBER_ROLES };

export const MEMBER_TABS = ["members", "invites"] as const;

export type MemberTab = (typeof MEMBER_TABS)[number];

// URL state for the members list. Mirrors ListMembersParams so the table talks
// to the server with whatever is in the query string (shareable, refresh-safe).
// `tab` switches Members vs Pending invites (owner/admin only).
export const membersSearchParams = {
	search: parseAsString.withDefault(""),
	role: parseAsStringLiteral(MEMBER_ROLES),
	page: parseAsInteger.withDefault(1),
	tab: parseAsStringLiteral(MEMBER_TABS).withDefault("members"),
};

// Shared defaults for discrete updates (role chips, pagination). Search uses
// these plus a debounce — see MembersToolbar.
export const membersSearchOptions = {
	history: "replace" as const,
	shallow: true,
};

export const membersSearchDebounce = debounce(300);
