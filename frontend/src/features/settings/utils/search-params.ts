import { parseAsStringLiteral } from "nuqs";

export const SETTINGS_SECTIONS = ["profile", "security", "sessions"] as const;

export type SettingsSection = (typeof SETTINGS_SECTIONS)[number];

export const SETTINGS_NAV: ReadonlyArray<{
	id: SettingsSection;
	label: string;
	description: string;
}> = [
	{ id: "profile", label: "Profile", description: "Name and email" },
	{ id: "security", label: "Security", description: "Password and 2FA" },
	{ id: "sessions", label: "Sessions", description: "Active logins" },
];

export const settingsSearchParams = {
	section: parseAsStringLiteral(SETTINGS_SECTIONS).withDefault("profile"),
};

export const settingsSearchOptions = {
	history: "replace" as const,
	shallow: true,
};
