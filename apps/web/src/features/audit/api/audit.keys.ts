import type { ListAuditParams } from "./audit.types.ts";

export const AUDIT_KEYS = {
	all: ["audit"] as const,
	workspace: (workspaceId: string) => [...AUDIT_KEYS.all, workspaceId] as const,
	list: (workspaceId: string, params: ListAuditParams = {}) =>
		[...AUDIT_KEYS.workspace(workspaceId), "list", params] as const,
};
