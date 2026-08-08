import { z } from "zod";

import { WORKSPACE_ICON_IDS } from "../lib/workspaceIcons.tsx";

export const workspaceSchema = z.object({
	name: z
		.string()
		.trim()
		.min(2, "Use at least 2 characters.")
		.max(50, "Keep it under 50 characters."),
	icon: z.enum(WORKSPACE_ICON_IDS),
});

export type WorkspaceValues = z.infer<typeof workspaceSchema>;
