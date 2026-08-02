import { z } from "zod";

import { ASSIGNABLE_ROLES } from "#/features/members/utils";

// Invite form — mirrors InviteInput (POST /members/invite). Owner is never
// assignable; role is constrained to admin|member.
export const inviteMemberSchema = z.object({
	email: z.email({ message: "Enter a valid email address." }),
	role: z.enum(ASSIGNABLE_ROLES, {
		message: "Choose a role.",
	}),
});

export type InviteMemberValues = z.infer<typeof inviteMemberSchema>;
