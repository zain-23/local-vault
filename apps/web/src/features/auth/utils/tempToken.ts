const TEMP_TOKEN_KEY = "lv_2fa_temp_token";

/** Short-lived login 2FA challenge token — sessionStorage so a tab refresh keeps it. */
export function setTempToken(token: string) {
	sessionStorage.setItem(TEMP_TOKEN_KEY, token);
}

export function getTempToken(): string | null {
	return sessionStorage.getItem(TEMP_TOKEN_KEY);
}

export function clearTempToken() {
	sessionStorage.removeItem(TEMP_TOKEN_KEY);
}
