export interface UpdateProfileInput {
	name?: string;
}

export interface AccountProfile {
	id: string;
	email: string;
	name: string;
	avatar_url?: string;
	onboarded: boolean;
	created_at: string;
}

export interface AccountSession {
	id: string;
	ip: string;
	user_agent: string;
	current: boolean;
	created_at: string;
	expires_at: string;
}
