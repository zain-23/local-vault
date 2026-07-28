package workspace

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/zain-23/local-vault/apps/server/internal/audit"
	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/id"
)

// Service holds workspace business logic
type Service struct {
	store *Store
	audit audit.Recorder
}

func NewService(store *Store, recorder audit.Recorder) *Service {
	return &Service{store: store, audit: recorder}
}

// uniqueSlug slugifies a name and appends -2, -3, ... until the slug is free
func (s *Service) uniqueSlug(ctx context.Context, name string) (string, error) {
	base := slugify(name)
	if base == "" {
		base = "workspace" // name was all symbols — fall back to a safe default
	}

	candidate := base
	for i := 2; ; i++ { // loop until we find a free slug
		exists, err := s.store.SlugExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// Create makes a workspace and its owner membership in one transaction
func (s *Service) Create(ctx context.Context, userID string, req CreateWorkspaceRequest) (*WorkspaceResponse, error) {
	slug, err := s.uniqueSlug(ctx, req.Name)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	now := time.Now()
	ws := &Workspace{
		ID:        id.Generate("ws_", 12),
		Name:      req.Name,
		Slug:      slug,
		OwnerID:   userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	mem := &Membership{
		ID:          id.Generate("mem_", 12),
		WorkspaceID: ws.ID,
		UserID:      userID,
		Role:        RoleOwner, // creator is always the owner
		JoinedAt:    now,
	}

	if err := s.store.CreateWorkspaceWithOwner(ctx, ws, mem); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: 	ws.ID,
		Action: 		"workspace.created",
		TargetType: 	"workspace",
		TargetID: 		ws.ID,
		TargetName: 	ws.Name,
	})

	return &WorkspaceResponse{Workspace: *ws, Role: RoleOwner}, nil
}

// List returns every workspace the user belongs to, each tagged with their role
func (s *Service) List(ctx context.Context, userID string) ([]WorkspaceResponse, error) {
	members, err := s.store.FindMembershipsByUserID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	// collect workspace ids + remember each role by workspace id
	ids := make([]string, 0, len(members))
	roleByWS := make(map[string]string, len(members))
	for _, m := range members {
		ids = append(ids, m.WorkspaceID)
		roleByWS[m.WorkspaceID] = m.Role
	}

	workspaces, err := s.store.FindWorkspacesByIDs(ctx, ids)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	// pair each workspace with the caller's role
	out := make([]WorkspaceResponse, 0, len(workspaces))
	for _, w := range workspaces {
		out = append(out, WorkspaceResponse{Workspace: w, Role: roleByWS[w.ID]})
	}
	return out, nil
}

// Get returns one workspace only if the caller is a member of it
func (s *Service) Get(ctx context.Context, userID, workspaceID string) (*WorkspaceResponse, error) {
	mem, err := s.store.GetMembership(ctx, workspaceID, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if mem == nil {
		return nil, apperror.New(404, "workspace not found") // non-members can't tell it exists
	}

	ws, err := s.store.FindWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if ws == nil {
		return nil, apperror.ErrNotFound // membership dangled (shouldn't happen) — treat as missing
	}
	return &WorkspaceResponse{Workspace: *ws, Role: mem.Role}, nil
}

// Update renames a workspace — RBAC (owner/admin) already enforced by middleware
func (s *Service) Update(ctx context.Context, workspaceID string, req UpdateWorkspaceRequest) (*Workspace, error) {
	err := s.store.UpdateWorkspace(ctx, workspaceID, bson.M{
		"name":       req.Name,
		"updated_at": time.Now(),
	})
	if err != nil {
		return nil, apperror.ErrInternal
	}

	ws, err := s.store.FindWorkspaceByID(ctx, workspaceID)
	if err != nil || ws == nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "workspace.updated",
		TargetType:  "workspace",
		TargetID:    workspaceID,
		TargetName:  ws.Name,
	})

	return ws, nil
}

// Delete removes a workspace + its memberships — RBAC (owner) already enforced by middleware
func (s *Service) Delete(ctx context.Context, workspaceID string) error {
	if err := s.store.DeleteWorkspaceCascade(ctx, workspaceID); err != nil {
		return apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "workspace.deleted",
		TargetType:  "workspace",
		TargetID:    workspaceID,
	})
	
	return nil
}