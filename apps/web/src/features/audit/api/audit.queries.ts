import { queryOptions } from "@tanstack/react-query";
import type { ListAuditParams } from "#/features/audit/api";
import { AUDIT_KEYS, auditService } from "#/features/audit/api";

export function auditEventsQuery(
	workspaceId: string,
	params: ListAuditParams = {},
) {
	return queryOptions({
		queryKey: AUDIT_KEYS.list(workspaceId, params),
		queryFn: async () => {
			const res = await auditService.list(workspaceId, params);
			return res.data;
		},
		staleTime: 20_000,
	});
}
