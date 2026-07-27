export interface UpdateProfileInput {
	name?: string;
	avatar_url?: string;
}

export interface AccountProfile {
	id: string;
	email: string;
	name: string;
	avatar_url?: string;
	two_factor_enabled: boolean;
	onboarded: boolean;
	created_at: string;
}

export interface ChangePasswordInput {
	current_password: string;
	new_password: string;
}

export interface Enable2FAResult {
	secret: string;
	otpauth_url: string;
}

export interface Verify2FAResult {
	backup_codes: string[];
}

export interface Disable2FAInput {
	totp_code?: string;
	backup_code?: string;
}

export interface AccountSession {
	id: string;
	ip: string;
	user_agent: string;
	current: boolean;
	created_at: string;
	expires_at: string;
}
