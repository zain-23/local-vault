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

  // Prefill/resume — the workspaces the caller belongs to. During onboarding
  // this is at most one; StepWorkspace prefills its input from the first.
  listWorkspaces() {
    return this.client.get<WorkspaceResponse[]>("/workspaces/");
  }

  // Step 1 — create the first workspace. Trailing slash matches the server's
  // route group (POST /api/v1/workspaces/). Returns the workspace + the
  // caller's role ("owner").
  createWorkspace(input: CreateWorkspaceInput) {
    return this.client.post<WorkspaceResponse>("/workspaces/", input);
  }

  // Step 1 (rename path) — PUT /workspaces/:id. Used when the user edits a
  // prefilled name instead of creating a fresh workspace. The server returns
  // the updated workspace bare (no role), so this resolves to Workspace.
  updateWorkspace(id: string, input: UpdateWorkspaceInput) {
    return this.client.put<Workspace>(`/workspaces/${id}`, input);
  }

  // Step 3 — list the caller's authorized devices. Polled while we wait for a
  // terminal to finish `lv login`; a non-empty list means one connected.
  listDevices() {
    return this.client.get<Device[]>("/device/");
  }

  // Final step — flip the account's `onboarded` flag. The server's UpdateProfile
  // only writes non-empty fields, so sending just this leaves name/avatar alone.
  completeOnboarding() {
    return this.client.put<null>("/account/me", {
      onboarded: true,
    } satisfies CompleteOnboardingInput);
  }
}

// Shared singleton for app use; construct with a custom client in tests.
export const onboardingService = new OnboardingService();
export { OnboardingService };
