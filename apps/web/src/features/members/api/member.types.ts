// Mirrors server/internal/member/dto.go + common/pagination. Snake_case matches
// the JSON the API returns. The UI's camelCase Member in types/index.ts stays
// until the page is wired — map at the boundary then.

export type MemberRole = "owner" | "admin" | "member";

// Roles an invite / role-change may assign. Owner is never assignable through
// these endpoints (server `oneof` rejects it).
export type AssignableRole = Exclude<MemberRole, "owner">;

// Pagination meta on every list response (server pagination.Meta).
export interface PageMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface Page<T> {
  items: T[];
  meta: PageMeta;
}

// One row of GET /workspaces/:wid/members.
export interface Member {
  user_id: string;
  name: string;
  email: string;
  avatar_url?: string;
  role: MemberRole;
  joined_at: string;
}

// One pending invite from GET /workspaces/:wid/members/invites. Never includes
// the raw token.
export interface Invite {
  id: string;
  email: string;
  role: AssignableRole;
  invited_by: string;
  created_at: string;
  expires_at: string;
}

// GET /members query params. Empty role/search means "no filter".
export interface ListMembersParams {
  page?: number;
  limit?: number;
  role?: MemberRole;
  search?: string;
}

// POST /members/invite body.
export interface InviteInput {
  email: string;
  role: AssignableRole;
}

// PUT /members/:userId/role body.
export interface ChangeRoleInput {
  role: AssignableRole;
}

// POST /members/join body — the raw token from the invite email link.
export interface JoinInput {
  token: string;
}

// POST /members/join response.
export interface JoinResult {
  workspace_id: string;
  role: AssignableRole;
}
