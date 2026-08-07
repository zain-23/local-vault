export interface User {
	id: string;
	email: string;
	name: string;
	oauth_provider?: string;
	avatar_url?: string;
	onboarded: boolean;
	created_at: string;
	updated_at: string;
}

export interface RefreshResponse {
	access_token: string;
}
