export const VAULT_KEYS = {
  all: ["vaults"] as const,
  workspace: (workspaceId: string) => [...VAULT_KEYS.all, workspaceId] as const,
  list: (workspaceId: string) =>
    [...VAULT_KEYS.workspace(workspaceId), "list"] as const,
};
