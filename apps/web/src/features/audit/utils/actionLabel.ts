/** Last segment of an action string — e.g. vault.collaborator.invited → invited */
export function actionVerb(action: string): string {
	const parts = action.split(".").filter(Boolean);
	return parts.at(-1) ?? action;
}

/** Short human label for the full action path. */
export function actionLabel(action: string): string {
	return actionVerb(action).replace(/_/g, " ");
}
