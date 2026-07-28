import { type ApiClient, api } from "#/services/api";
import type {
  ChangeRoleInput,
  Invite,
  InviteInput,
  JoinInput,
  JoinResult,
  ListMembersParams,
  Member,
  Page,
} from "./member.types.ts";

// One method per route under /api/v1/workspaces/:wid/members
// (server/internal/member/routes.go). Generics are the envelope `data` type.
// The client is injected so tests can pass a stub.
class MemberService {
  constructor(private readonly client: ApiClient = api) {}

  private base(workspaceId: string) {
    return `/workspaces/${encodeURIComponent(workspaceId)}/members`;
  }

  // GET / — paginated roster. Optional role/search filters.
  listMembers(workspaceId: string, params: ListMembersParams = {}) {
    return this.client.get<Page<Member>>(`${this.base(workspaceId)}`, {
      params: {
        page: params.page,
        limit: params.limit,
        role: params.role,
        search: params.search,
      },
    });
  }

  // POST /invite — create a pending invite and enqueue the email.
  invite(workspaceId: string, input: InviteInput) {
    return this.client.post<Invite>(`${this.base(workspaceId)}/invite`, input);
  }

  // GET /invites — pending invites (never includes the raw token).
  listInvites(workspaceId: string) {
    return this.client.get<Invite[]>(`${this.base(workspaceId)}/invites`);
  }

  // DELETE /invites/:id — cancel a pending invite.
  cancelInvite(workspaceId: string, inviteId: string) {
    return this.client.delete<null>(
      `${this.base(workspaceId)}/invites/${encodeURIComponent(inviteId)}`,
    );
  }

  // PUT /:userId/role — change a member's role (admin|member only).
  changeRole(workspaceId: string, userId: string, input: ChangeRoleInput) {
    return this.client.put<Member>(
      `${this.base(workspaceId)}/${encodeURIComponent(userId)}/role`,
      input,
    );
  }

  // DELETE /:userId — remove a member from the workspace.
  removeMember(workspaceId: string, userId: string) {
    return this.client.delete<null>(
      `${this.base(workspaceId)}/${encodeURIComponent(userId)}`,
    );
  }

  // POST /join — accept an invite token. Auth required; no RBAC (joiner isn't
  // a member yet). workspaceId in the path must match the invite's workspace.
  join(workspaceId: string, input: JoinInput) {
    return this.client.post<JoinResult>(
      `${this.base(workspaceId)}/join`,
      input,
    );
  }
}

// Shared singleton for app use; construct with a custom client in tests.
export const memberService = new MemberService();
export { MemberService };
