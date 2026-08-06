import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

import type { DeviceDecisionAction } from "#/features/device/api";
import { DEVICE_KEYS, deviceService } from "#/features/device/api";
import type { ApiResponse } from "#/services/api";

// Approve or deny the pending request. The action is the mutation variable, so
// the screen calls decide.mutate("approve" | "deny") and reads decide.variables
// to render the matching outcome. On success we refresh the details query so a
// re-read reflects the new status. Errors surface as a toast.
export function useDecideDevice(userCode: string) {
	return useMutation<ApiResponse<null>, Error, DeviceDecisionAction>({
		mutationKey: DEVICE_KEYS.decide(userCode),
		mutationFn: (action) => deviceService.decide(userCode, { action }),
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
