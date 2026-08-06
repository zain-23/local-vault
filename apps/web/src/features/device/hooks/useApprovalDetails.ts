import { useQuery } from "@tanstack/react-query";

import { approvalDetailsQuery } from "#/features/device/api";

// Loads the request behind a user code so the screen can show which device
// (name + IP) is asking before the user approves or denies. Disabled when the
// URL carries no code — there's nothing to look up, and the route sends
// unauthenticated visitors to login before this ever runs.
export function useApprovalDetails(userCode: string | undefined) {
	return useQuery({
		...approvalDetailsQuery(userCode ?? ""),
		enabled: !!userCode,
	});
}
