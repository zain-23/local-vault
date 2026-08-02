import { type ApiClient, api } from "#/services/api";
import type {
	AuthTokens,
	LoginInput,
	LoginResult,
	RefreshResponse,
	ResetPasswordInput,
	SignupInput,
	User,
} from "./auth.types.ts";

// One method per server route (server/internal/auth/routes.go). Generics are the
// endpoint's `data` type; the useful text for message-only routes lives in the
// returned envelope's `message`. Client is injected for easy test mocking.
class AuthService {
	constructor(private readonly client: ApiClient = api) {}

	// Message-only: data is "" — read the envelope's `message`.
	signup(input: SignupInput) {
		return this.client.post<string>("/auth/signup", input);
	}

	// Resolves to the signed-in User (cookies set) OR a 2FA challenge.
	login(input: LoginInput) {
		return this.client.post<LoginResult>("/auth/login", input);
	}

	// Refresh reads the refresh_token cookie server-side; no body needed.
	refresh() {
		return this.client.post<RefreshResponse>("/auth/refresh");
	}

	logout() {
		return this.client.post<null>("/auth/logout");
	}

	// Message-only: token rides in the query, server verifies and returns no body.
	verifyEmail(token: string) {
		return this.client.post<null>("/auth/verify-email", undefined, {
			params: { token },
		});
	}

	forgotPassword(email: string) {
		return this.client.post<null>("/auth/forgot-password", { email });
	}

	resetPassword(input: ResetPasswordInput) {
		return this.client.post<null>("/auth/reset-password", input);
	}

	sendMagicLink(email: string) {
		return this.client.post<null>("/auth/magic-link", { email });
	}

	// Consuming the magic-link token logs the user in — token bundle in the body.
	verifyMagicLink(token: string) {
		return this.client.post<AuthTokens>("/auth/magic-link/verify", { token });
	}

	// OAuth login
	oauthUrl(provider: "google" = "google") {
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
