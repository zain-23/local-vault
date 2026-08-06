import { queryOptions } from "@tanstack/react-query";

import { DEVICE_KEYS, deviceService } from "#/features/device/api";

// The pending request behind a user code, shown on the approval screen. Keyed by
// code so distinct codes never collide. retry:false — a 404 (unknown/expired
// code) is a real answer, not a transient failure worth retrying. staleTime:0 —
// we always want the current status when the screen (re)mounts.
export function approvalDetailsQuery(userCode: string) {
	return queryOptions({
		queryKey: DEVICE_KEYS.approval(userCode),
		queryFn: async () => {
			const res = await deviceService.approvalDetails(userCode);
			return res.data;
		},
		retry: false,
		staleTime: 0,
	});
}
