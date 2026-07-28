export interface User {
	id: string;
	email: string;
	name: string;
	oauth_provider?: string;
	avatar_url?: string;
	email_verified: boolean;
	two_factor_enabled: boolean;
	onboarded: boolean;
	created_at: string;
	updated_at: string;
}

export interface AuthTokens {
	access_token: string;
	refresh_token: string;
	user: User;
}

// The 2FA branch of login: account has 2FA on, so no session is issued yet.
export interface Login2FARequired {
	requires_2fa: true;
	temp_token: string;
}

// POST /auth/login resolves to EITHER the signed-in User (cookies set) OR a 2FA
// challenge — never a combined object. Narrow with `isTwoFactorRequired`.
export type LoginResult = User | Login2FARequired;

export function isTwoFactorRequired(
	result: LoginResult,
): result is Login2FARequired {
	return "requires_2fa" in result;
}

export interface RefreshResponse {
	access_token: string;
}

// ---- request payloads ----
export interface SignupInput {
	email: string;
	password: string;
	name: string;
}

export interface LoginInput {
	email: string;
	password: string;
}

export interface Verify2FAInput {
	temp_token: string;
	totp_code: string;
}

export interface ResetPasswordInput {
	token: string;
	new_password: string;
}
