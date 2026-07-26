package member

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zain-23/local-vault/server/internal/audit"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/id"
	"github.com/zain-23/local-vault/server/internal/common/pagination"
	"github.com/zain-23/local-vault/server/internal/config"
	"github.com/zain-23/local-vault/server/internal/email"
)

// Service holds member business logic.
type Service struct {
	store     *Store
	publisher *email.Publisher // enqueues invite emails onto RabbitMQ
	cfg       config.Config    // for FrontendURL when building the join link
	audit 	  audit.Recorder
}

func NewService(store *Store, pub *email.Publisher, cfg config.Config, recorder audit.Recorder) *Service {
	return &Service{store: store, publisher: pub, cfg: cfg, audit: recorder}
}

// enrich joins membership rows with user display fields in one users query.
func (s *Service) enrich(ctx context.Context, mems []Membership) ([]MemberResponse, error) {
	ids := make([]string, 0, len(mems))
	for _, m := range mems {
		ids = append(ids, m.UserID)
	}

	users, err := s.store.FindUsersByIDs(ctx, ids)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	// index users by id so each membership can find its display fields
	byID := make(map[string]UserInfo, len(users))
	for _, u := range users {
		byID[u.ID] = u
	}

	out := make([]MemberResponse, 0, len(mems))
	for _, m := range mems {
		u := byID[m.UserID] // zero-value UserInfo if the user doc is missing (shouldn't happen)
		out = append(out, MemberResponse{
			UserID:    m.UserID,
			Name:      u.Name,
			Email:     u.Email,
			AvatarURL: u.AvatarURL,
			Role:      m.Role,
			JoinedAt:  m.JoinedAt,
		})
	}
	return out, nil
}

// List returns one page of a workspace's members (with display fields), filtered
// by role/search, plus the pagination meta. The store's aggregation does the
// user join, so enrich() isn't needed on this path.
func (s *Service) List(ctx context.Context, workspaceID string, q ListMembersQuery) (*pagination.Page[MemberResponse], error) {
	items, total, err := s.store.ListMembersPaginated(ctx, workspaceID, q)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return &pagination.Page[MemberResponse]{
		Items: items,
		Meta:  pagination.NewMeta(q.Params, total),
	}, nil
}

// Invite creates a pending invite and enqueues the invite email.
// Role is already validated (admin|member) by the DTO before this runs.
func (s *Service) Invite(ctx context.Context, workspaceID, invitedBy string, req InviteRequest) (*InviteResponse, error) {
	// NOTE: name this local `addr`, not `email` — a local named `email` would
	// shadow the imported `email` package and break `email.EmailJob` below.
	addr := normalizeEmail(req.Email)

	// Reject if this email already belongs to a member of this workspace.
	user, err := s.store.FindUserByEmail(ctx, addr)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if user != nil {
		isMember, err := s.store.MembershipExists(ctx, workspaceID, user.ID)
		if err != nil {
			return nil, apperror.ErrInternal
		}
		if isMember {
			return nil, apperror.New(409, "user is already a member")
		}
	}

	// Reject if a live pending invite already exists for this email here.
	pending, err := s.store.PendingInviteExists(ctx, workspaceID, addr)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if pending {
		return nil, apperror.New(409, "an invite is already pending for this email")
	}

	// Generate the token; store only its hash. The raw token goes in the email link.
	rawToken, err := generateInviteToken()
	if err != nil {
		return nil, apperror.ErrInternal
	}

	now := time.Now()
	inv := &Invite{
		ID:          id.Generate("inv_", 12),
		WorkspaceID: workspaceID,
		Email:       addr,
		Role:        req.Role,
		InvitedBy:   invitedBy,
		TokenHash:   sha256Hex(rawToken),
		Status:      StatusPending,
		CreatedAt:   now,
		ExpiresAt:   now.Add(7 * 24 * time.Hour), // 7-day expiry
	}
	if err := s.store.CreateInvite(ctx, inv); err != nil {
		return nil, apperror.ErrInternal
	}

	// Look up the workspace name for the email; fall back gracefully if missing.
	wsName, err := s.store.FindWorkspaceName(ctx, workspaceID)
	if err != nil || wsName == "" {
		wsName = "a workspace"
	}

	// Enqueue the invite email — a queue hiccup must not fail the request.
	// workspaceId is required by POST /members/join (:wid) — the token alone
	// isn't enough for the frontend to build the request path.
	joinURL := fmt.Sprintf("%s/workspaces/join?token=%s&workspaceId=%s", s.cfg.FrontendURL, rawToken, workspaceID)
	job := email.EmailJob{
		Kind: email.KindWorkspaceInvite,
		To:   addr,
		Name: wsName, // the invite template greets with the workspace name
		URL:  joinURL,
	}
	if err := s.publisher.Publish(ctx, job); err != nil {
		log.Printf("⚠️ failed to enqueue invite email for %s: %v", addr, err)
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "member.invited",
		TargetType:  "invite",
		TargetID:    inv.ID,
		TargetName:  inv.Email,
		Details:     map[string]any{"role": inv.Role},
	})

	return &InviteResponse{
		ID:        inv.ID,
		Email:     inv.Email,
		Role:      inv.Role,
		InvitedBy: inv.InvitedBy,
		CreatedAt: inv.CreatedAt,
		ExpiresAt: inv.ExpiresAt,
	}, nil
}

// ListInvites returns pending invites for a workspace (never the token).
func (s *Service) ListInvites(ctx context.Context, workspaceID string) ([]InviteResponse, error) {
	invites, err := s.store.ListPendingInvites(ctx, workspaceID)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	out := make([]InviteResponse, 0, len(invites))
	for _, inv := range invites {
		out = append(out, InviteResponse{
			ID:        inv.ID,
			Email:     inv.Email,
			Role:      inv.Role,
			InvitedBy: inv.InvitedBy,
			CreatedAt: inv.CreatedAt,
			ExpiresAt: inv.ExpiresAt,
		})
	}
	return out, nil
}

// CancelInvite deletes a pending invite. The invite must belong to this workspace,
// otherwise we 404 (never confirm another workspace's invite exists).
func (s *Service) CancelInvite(ctx context.Context, workspaceID, inviteID string) error {
	inv, err := s.store.FindInviteByID(ctx, inviteID)
	if err != nil {
		return apperror.ErrInternal
	}
	if inv == nil || inv.WorkspaceID != workspaceID {
		return apperror.New(404, "invite not found")
	}
	if err := s.store.DeleteInvite(ctx, inviteID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}

// Join accepts an invite for the authenticated caller. The caller's email must
// match the invite's email. Creates the membership + marks the invite accepted atomically.
func (s *Service) Join(ctx context.Context, workspaceID, userID, userEmail string, req JoinRequest) (*JoinResponse, error) {
	inv, err := s.store.FindInviteByTokenHash(ctx, sha256Hex(req.Token))
	if err != nil {
		return nil, apperror.ErrInternal
	}
	// Missing, already used, or expired all look the same to the caller.
	if inv == nil || inv.Status != StatusPending || inv.ExpiresAt.Before(time.Now()) {
		return nil, apperror.New(400, "invalid or expired invite")
	}
	// The token is workspace-scoped — the :wid in the path must match.
	if inv.WorkspaceID != workspaceID {
		return nil, apperror.New(400, "invite does not match this workspace")
	}
	// The invite is bound to a specific email; only that person may accept it.
	if normalizeEmail(userEmail) != inv.Email {
		return nil, apperror.New(403, "this invite was sent to a different email")
	}

	// Guard against a double-join (unique index also enforces this at the DB level).
	isMember, err := s.store.MembershipExists(ctx, workspaceID, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if isMember {
		return nil, apperror.New(409, "you are already a member")
	}

	now := time.Now()
	mem := &Membership{
		ID:          id.Generate("mem_", 12),
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        inv.Role, // role comes from the invite (admin|member)
		InvitedBy:   inv.InvitedBy,
		JoinedAt:    now,
	}
	if err := s.store.AcceptInvite(ctx, inv, mem); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "member.joined",
		TargetType:  "user",
		TargetID:    userID,
		Details:     map[string]any{"role": inv.Role},
	})

	return &JoinResponse{WorkspaceID: workspaceID, Role: inv.Role}, nil
}

// ChangeRole changes a member's role. The owner's role can never be changed here.
// Role is already validated (admin|member) by the DTO.
func (s *Service) ChangeRole(ctx context.Context, workspaceID, targetUserID string, req ChangeRoleRequest) (*MemberResponse, error) {
	target, err := s.store.GetMembership(ctx, workspaceID, targetUserID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if target == nil {
		return nil, apperror.New(404, "member not found")
	}
	if target.Role == RoleOwner {
		return nil, apperror.New(403, "cannot change the owner's role")
	}

	if err := s.store.UpdateMemberRole(ctx, workspaceID, targetUserID, req.Role); err != nil {
		return nil, apperror.ErrInternal
	}

	// Return the updated member, enriched with display fields (reuse the List helper).
	target.Role = req.Role
	enriched, err := s.enrich(ctx, []Membership{*target})
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "member.role.changed",
		TargetType:  "user",
		TargetID:    targetUserID,
		Details:     map[string]any{"new_role": req.Role},
	})

	return &enriched[0], nil
}

// RemoveMember removes a member. The owner can never be removed here.
func (s *Service) RemoveMember(ctx context.Context, workspaceID, targetUserID string) error {
	target, err := s.store.GetMembership(ctx, workspaceID, targetUserID)
	if err != nil {
		return apperror.ErrInternal
	}
	if target == nil {
		return apperror.New(404, "member not found")
	}
	if target.Role == RoleOwner {
		return apperror.New(403, "cannot remove the workspace owner")
	}
	if err := s.store.DeleteMembership(ctx, workspaceID, targetUserID); err != nil {
		return apperror.ErrInternal
	}
	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "member.removed",
		TargetType:  "user",
		TargetID:    targetUserID,
	})
	return nil
}
