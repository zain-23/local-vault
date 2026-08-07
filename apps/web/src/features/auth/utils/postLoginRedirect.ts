const POST_LOGIN_REDIRECT_KEY = "lv_post_login_redirect";

/**
 * Carries a guard's `?redirect=` across the GitHub OAuth round trip. The
 * backend's OAuth callback always lands on /dashboard (it has no way to pass
 * the original redirect through the provider), so the target is stashed here
 * before leaving the app and consumed once on the other side.
 */
export function setPostLoginRedirect(target: string) {
	sessionStorage.setItem(POST_LOGIN_REDIRECT_KEY, target);
}

/** Reads and clears the stashed target in one step — it's meant to fire once. */
export function consumePostLoginRedirect(): string | null {
	const target = sessionStorage.getItem(POST_LOGIN_REDIRECT_KEY);
	sessionStorage.removeItem(POST_LOGIN_REDIRECT_KEY);
	return target;
}
