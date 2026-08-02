import { useQuery } from "@tanstack/react-query";

import { devicesQuery } from "#/features/onboarding/api";

// How often to re-check for a linked terminal while Step 3 is open.
const POLL_INTERVAL_MS = 3_000;

// Step 3 — waits for a terminal to finish `lv login`. Polls the device list and
// resolves `connected` once at least one authorized device shows up. Polling
// stops the moment we're connected so we don't hammer the API after success.
export function useDeviceConnection() {
	const query = useQuery({
		...devicesQuery,
		refetchInterval: (q) => {
			const devices = q.state.data ?? [];
			return devices.length > 0 ? false : POLL_INTERVAL_MS;
		},
	});

	const devices = query.data ?? [];

	return {
		devices,
		connected: devices.length > 0,
		isWaiting: query.isPending || (devices.length === 0 && query.isFetching),
		error: query.error,
	};
}
