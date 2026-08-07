// Keys mirror the ?error= values the OAuth callback redirects with —
// see apps/server/internal/auth/oauth.go's HandleCallback. Keep both in sync.
const AUTH_ERROR_MESSAGES: Record<string, string> = {
	"invalid state": "That login link expired. Please try again.",
	"oauth failed": "GitHub sign-in failed. Please try again.",
	"profile failed": "Couldn't read your GitHub profile. Please try again.",
	"account failed": "Couldn't sign you in. Please try again.",
	"session failed": "Couldn't start your session. Please try again.",
};

export function authErrorMessage(code: string): string {
	return (
		AUTH_ERROR_MESSAGES[code] ??
		"Something went wrong signing in. Please try again."
	);
}
