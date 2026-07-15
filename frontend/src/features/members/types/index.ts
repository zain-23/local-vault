// Mirrors the server's MemberResponse (server/internal/member/dto.go), camelCased
// for the frontend. Swap the mock data source for the real API later and this
// type stays the same.
export type MemberRole = "owner" | "admin" | "member";

export interface Member {
	userId: string;
	name: string;
	email: string;
	avatarUrl?: string;
	role: MemberRole;
	joinedAt: string; // ISO timestamp
}

// The role values a member can be invited as / changed to. The owner role is
// never assignable through the UI (server rejects it too).
export const ASSIGNABLE_ROLES: Exclude<MemberRole, "owner">[] = [
	"admin",
	"member",
];
