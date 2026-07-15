import { create } from "zustand";

// The currently active workspace. The member API is workspace-scoped
// (/workspaces/:wid/members), so this id is what feature hooks will read to
// build their requests. The sidebar's WorkspaceSwitcher sets it.
export interface Workspace {
	id: string;
	name: string;
	plan: string;
}

interface WorkspaceState {
	active: Workspace | null;
	setActive: (workspace: Workspace) => void;
}

// Seed with a placeholder until real workspaces are wired (onboarding creates one).
const PLACEHOLDER_WORKSPACE: Workspace = {
	id: "ws_placeholder",
	name: "Kodexo Labs",
	plan: "Free",
};

export const useWorkspaceStore = create<WorkspaceState>((set) => ({
	active: PLACEHOLDER_WORKSPACE,
	setActive: (workspace) => set({ active: workspace }),
}));
