import { type ApiClient, api } from "#/services/api";

import type {
	AccountProfile,
	AccountSession,
	UpdateProfileInput,
} from "./settings.types.ts";

// One method per route under /api/v1/account. The client is injectable so this
// boundary remains straightforward to test with a stubbed API client.
class SettingsService {
	constructor(private readonly client: ApiClient = api) {}

	updateProfile(input: UpdateProfileInput) {
		return this.client.put<AccountProfile>("/account/me", input);
	}

	listSessions() {
		return this.client.get<AccountSession[]>("/account/sessions");
	}

	revokeSession(sessionId: string) {
		return this.client.delete<null>(
			`/account/sessions/${encodeURIComponent(sessionId)}`,
		);
	}

	revokeOtherSessions() {
		return this.client.delete<null>("/account/sessions");
	}
}

export const settingsService = new SettingsService();
export { SettingsService };
