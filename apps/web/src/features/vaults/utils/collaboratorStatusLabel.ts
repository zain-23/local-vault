export function collaboratorStatusLabel(status: string): string {
	switch (status) {
		case "pending":
			return "Pending";
		case "active":
			return "Active";
		case "revoked":
			return "Revoked";
		default:
			return status;
	}
}
