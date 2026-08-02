import { type ApiClient, api } from "#/services/api";
import type {
	CompleteOnboardingInput,
	CreateWorkspaceInput,
	Device,
	UpdateWorkspaceInput,
	Workspace,
	WorkspaceResponse,
} from "./onboarding.types.ts";

// One method per server route the onboarding wizard touches. It spans three
// domains (workspaces, device, account) on purpose: these calls only exist to
// drive onboarding, so keeping them together keeps the flow readable in one
// place. Generics are the endpoint's envelope `data` type. The client is
// injected so tests can pass a stub.
class OnboardingService {
	constructor(private readonly client: ApiClient = api) {}

	listWorkspaces() {
		return this.client.get<WorkspaceResponse[]>("/workspaces");
	}

	createWorkspace(input: CreateWorkspaceInput) {
		return this.client.post<WorkspaceResponse>("/workspaces", input);
	}

	updateWorkspace(id: string, input: UpdateWorkspaceInput) {
		return this.client.put<Workspace>(`/workspaces/${id}`, input);
	}

	listDevices() {
		return this.client.get<Device[]>("/device");
	}

	completeOnboarding() {
		return this.client.put<null>("/account/me", {
			onboarded: true,
		} satisfies CompleteOnboardingInput);
	}
}

// Shared singleton for app use; construct with a custom client in tests.
export const onboardingService = new OnboardingService();
export { OnboardingService };
