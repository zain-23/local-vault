import { type ApiClient, api } from "#/services/api";

import type {
	AccountProfile,
	AccountSession,
	ChangePasswordInput,
	Disable2FAInput,
	Enable2FAResult,
	UpdateProfileInput,
	Verify2FAResult,
} from "./settings.types.ts";

// One method per route under /api/v1/account. The client is injectable so this
// boundary remains straightforward to test with a stubbed API client.
class SettingsService {
	constructor(private readonly client: ApiClient = api) {}

	updateProfile(input: UpdateProfileInput) {
		return this.client.put<AccountProfile>("/account/me", input);
	}

	changePassword(input: ChangePasswordInput) {
		return this.client.put<null>("/account/password", input);
	}

	enable2FA() {
		return this.client.post<Enable2FAResult>("/account/2fa/enable");
	}

	verify2FA(totpCode: string) {
		return this.client.post<Verify2FAResult>("/account/2fa/verify", {
			totp_code: totpCode,
		});
	}

	disable2FA(input: Disable2FAInput) {
		return this.client.post<null>("/account/2fa/disable", input);
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
