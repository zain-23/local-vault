import { create } from "zustand";

import type { MemberRole } from "#/features/members/types";

// The currently active workspace. The member API is workspace-scoped
// (/workspaces/:wid/members), so this id is what feature hooks will read to
// build their requests. `role` is the caller's membership in this workspace
// (from GET /workspaces) — used to gate invite UI.
export interface Workspace {
	id: string;
	name: string;
	plan: string;
	role: MemberRole | null;
}

interface WorkspaceState {
	active: Workspace | null;
	setActive: (workspace: Workspace) => void;
}

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
	active: null,
	setActive: (workspace) => set({ active: workspace }),
}));
