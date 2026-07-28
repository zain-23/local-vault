import { useQuery } from "@tanstack/react-query";
import type { ListAuditParams } from "#/features/audit/api";
import { auditEventsQuery } from "#/features/audit/api";
import { useWorkspaceStore } from "#/stores";

export function useAuditEvents(params: ListAuditParams = {}) {
  const workspaceId = useWorkspaceStore((s) => s.active?.id);
  return useQuery({
    ...auditEventsQuery(workspaceId ?? "", params),
    enabled: !!workspaceId,
  });
}
