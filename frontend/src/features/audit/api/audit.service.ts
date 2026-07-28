import { type ApiClient, api } from "#/services/api";

import type { AuditEvent, ListAuditParams, Page } from "./audit.types.ts";

class AuditService {
  constructor(private readonly client: ApiClient = api) {}

  private base(workspaceId: string) {
    return `/workspaces/${encodeURIComponent(workspaceId)}/audit`;
  }

  // GET /workspaces/:wid/audit — paginated, filterable log (no export).
  list(workspaceId: string, params: ListAuditParams = {}) {
    return this.client.get<Page<AuditEvent>>(this.base(workspaceId), {
      params: {
        page: params.page,
        limit: params.limit,
        action: params.action,
        action_prefix: params.action_prefix,
        actor_id: params.actor_id,
        device_id: params.device_id,
        from: params.from,
        to: params.to,
      },
    });
  }
}

export const auditService = new AuditService();
export { AuditService };
