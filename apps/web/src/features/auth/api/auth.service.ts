import { type ApiClient, api } from "#/services/api";
import type { RefreshResponse, User } from "./auth.types.ts";

// One method per server route (server/internal/auth/routes.go). Client is
// injected for easy test mocking.
class AuthService {
	constructor(private readonly client: ApiClient = api) {}

	// Refresh reads the refresh_token cookie server-side; no body needed.
	refresh() {
		return this.client.post<RefreshResponse>("/auth/refresh");
	}

	logout() {
		return this.client.post<null>("/auth/logout");
	}

	// GitHub is the only OAuth provider — this is a full-page redirect, not a fetch.
	oauthUrl(provider: "github" = "github") {
		return `${import.meta.env.VITE_API_URL}/auth/oauth/${provider}`;
	}

	// Resolve to the signed-in User; 401 when there is no valid session
	me() {
		return this.client.get<User>("/account/me");
	}
}

// Shared singleton for app use; construct with a custom client in tests.
export const authService = new AuthService();
export { AuthService };
