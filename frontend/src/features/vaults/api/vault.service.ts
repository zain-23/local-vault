import { type ApiClient, api } from "#/services/api";
import type { VaultSummary } from "./vault.types.ts";

class VaultService {
  constructor(private readonly client: ApiClient = api) {}

  private base(workspaceId: string) {
    return `/workspaces/${encodeURIComponent(workspaceId)}/vaults`;
  }

  list(workspaceId: string) {
    return this.client.get<VaultSummary[]>(`${this.base(workspaceId)}`);
  }
}

export const vaultService = new VaultService();
