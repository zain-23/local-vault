import { type ApiClient, api } from "#/services/api";
import type { ApprovalDetails, DecisionInput, Device } from "./device.types.ts";

// One method per authenticated device route (server/internal/device/routes.go).
// The public CLI routes (authorize/poll) are the terminal's job, not the
// browser's, so they're intentionally absent here. Generics are the endpoint's
// envelope `data` type; the client is injected so tests can pass a stub.
class DeviceService {
	constructor(private readonly client: ApiClient = api) {}

	// GET /device/authorize/:userCode — the request behind a user code. `userCode`
	// is the hyphenated form the server stores and the URL carries (WDJF-X4K2);
	// it's encoded in case that ever changes. Auth required.
	approvalDetails(userCode: string) {
		return this.client.get<ApprovalDetails>(
			`/device/authorize/${encodeURIComponent(userCode)}`,
		);
	}

	// PUT /device/authorize/:userCode — approve or deny. Message-only (data null);
	// the useful text is the envelope's `message`.
	decide(userCode: string, input: DecisionInput) {
		return this.client.put<null>(
			`/device/authorize/${encodeURIComponent(userCode)}`,
			input,
		);
	}

	// GET /device — the caller's authorized devices. For the management UI (future).
	listDevices() {
		return this.client.get<Device[]>("/device/");
	}

	// DELETE /device/:id — revoke a device. For the management UI (future).
	revokeDevice(id: string) {
		return this.client.delete<null>(`/device/${encodeURIComponent(id)}`);
	}
}

// Shared singleton for app use; construct with a custom client in tests.
export const deviceService = new DeviceService();
export { DeviceService };
