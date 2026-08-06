import { z } from "zod";

export const workspaceSchema = z.object({
	name: z
		.string()
		.trim()
		.min(2, "Use at least 2 characters.")
		.max(50, "Keep it under 50 characters."),
});

export type WorkspaceValues = z.infer<typeof workspaceSchema>;
