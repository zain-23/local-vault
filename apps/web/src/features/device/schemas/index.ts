import { z } from "zod";

// The server's user-code alphabet (server/internal/device/code.go): digits 2-9
// and A-Z minus the ambiguous I, L, O, 0, 1. Kept in sync so the browser rejects
// codes the server could never have issued.
export const USER_CODE_ALPHABET = "23456789ABCDEFGHJKMNPQRSTUVWXYZ";
export const USER_CODE_LENGTH = 8; // 4 + 4, joined by a hyphen for display

const alphabetClass = `[${USER_CODE_ALPHABET}]`;

// Validates the raw 8-character code (no hyphen) — the input strips the dash so
// paste and typing share one value.
export const deviceCodeSchema = z.object({
	code: z
		.string()
		.length(USER_CODE_LENGTH, "Enter all 8 characters.")
		.regex(
			new RegExp(`^${alphabetClass}{${USER_CODE_LENGTH}}$`),
			"That code contains characters we don't use.",
		),
});

export type DeviceCodeValues = z.infer<typeof deviceCodeSchema>;

// Keeps only allowed characters and upper-cases them — used to sanitise both
// typed input and pasted values before they reach the form.
export function normalizeUserCode(raw: string): string {
	const allowed = new Set(USER_CODE_ALPHABET);
	return raw
		.toUpperCase()
		.split("")
		.filter((ch) => allowed.has(ch))
		.join("")
		.slice(0, USER_CODE_LENGTH);
}

// Turns the raw 8-char value into the hyphenated form the server stores and the
// route uses (WDJFX4K2 -> WDJF-X4K2).
export function toUserCodeParam(code: string): string {
	const c = normalizeUserCode(code);
	return `${c.slice(0, 4)}-${c.slice(4)}`;
}
