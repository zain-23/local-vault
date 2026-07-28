package vault

import (
	"context"
	"time"

	"github.com/zain-23/local-vault/server/internal/audit"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/id"
)

// Service holds the vault domain's business logic.
type Service struct {
	Store *Store
	audit audit.Recorder
}

func NewService(store *Store, recorder audit.Recorder) *Service {
	return &Service{Store: store, audit: recorder}
}

// getScoped loads a vault and enforces it belongs to workspaceID, else 404.
// Shared by every by-id operation so the scoping rule lives in one place.
func (s *Service) getScoped(ctx context.Context, workspaceID, vaultID string) (*Vault, error) {
	v, err := s.Store.FindVaultByID(ctx, vaultID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if v == nil || v.WorkspaceID != workspaceID {
		return nil, apperror.New(404, "vault not found") // never confirm cross-workspace ids
	}
	return v, nil
}


// Create registers a vault with the caller's device as the first peer.
func (s *Service) Create(ctx context.Context, workspaceID, createdBy string, req CreateVaultRequest) (*CreateVaultResponse, error) {
	now := time.Now()
	v := &Vault{
		ID:            id.Generate("vlt_", 12),
		WorkspaceID:   workspaceID,
		Name:          req.Name,
		CreatedBy:     createdBy,             // usr_ from the JWT
		OwnerDeviceID: req.OwnerDeviceID,     // P2P id of the creator's device
		Peers: []Peer{{ // seed the owner as peer #1
			DeviceID:        req.OwnerDeviceID,
			DeviceName:      req.OwnerName,
			PublicKey:       req.PublicKey,
			X25519PublicKey: req.X25519PublicKey,
			JoinedAt:        now,
		}},
		Tokens:    []Token{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.Store.CreateVault(ctx, v); err != nil {
		return nil, apperror.ErrInternal
	}
	
	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.created",
		TargetType:  "vault",
		TargetID:    v.ID,
		TargetName:  v.Name,
	})

	return &CreateVaultResponse{VaultID: v.ID, CreatedAt: v.CreatedAt}, nil
}


// List returns every vault in the workspace as lightweight summaries.
func (s *Service) List(ctx context.Context, workspaceID string) ([]VaultSummary, error) {
	vaults, err := s.Store.ListVaultsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	out := make([]VaultSummary, 0, len(vaults)) // len 0 → JSON [] not null
	for _, v := range vaults {
		out = append(out, VaultSummary{
			ID:            v.ID,
			Name:          v.Name,
			OwnerDeviceID: v.OwnerDeviceID,
			PeerCount:     len(v.Peers),
			HasSnapshot:   len(v.Snapshot) > 0,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		})
	}
	return out, nil
}

// Get returns the detail view (with peers) for one vault.
func (s *Service) Get(ctx context.Context, workspaceID, vaultID string) (*VaultResponse, error) {
	v, err := s.getScoped(ctx, workspaceID, vaultID)
	if err != nil {
		return nil, err
	}
	return &VaultResponse{
		ID:            v.ID,
		Name:          v.Name,
		WorkspaceID:   v.WorkspaceID,
		CreatedBy:     v.CreatedBy,
		OwnerDeviceID: v.OwnerDeviceID,
		Peers:         v.Peers,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
	}, nil
}

// Delete removes a vault (route restricts this to owner/admin).
func (s *Service) Delete(ctx context.Context, workspaceID, vaultID string) error {
	if _, err := s.getScoped(ctx, workspaceID, vaultID); err != nil {
		return err
	}
	if err := s.Store.DeleteVault(ctx, vaultID); err != nil {
		return apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.deleted",
		TargetType:  "vault",
		TargetID:    vaultID,
	})

	return nil
}

// PushSnapshot stores a new encrypted snapshot. deviceID is the caller's P2P id;
// only an existing peer may push.
func (s *Service) PushSnapshot(ctx context.Context, workspaceID, vaultID, deviceID string, snapshot []byte) (*SnapshotResponse, error) {
	v, err := s.getScoped(ctx, workspaceID, vaultID)
	if err != nil {
		return nil, err
	}
	if !hasPeer(v, deviceID) {
		return nil, apperror.New(403, "device is not a vault peer")
	}
	if err := s.Store.SetSnapshot(ctx, vaultID, snapshot); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.push",
		TargetType:  "vault",
		TargetID:    vaultID,
		Details:     map[string]any{"device_id": deviceID, "bytes": len(snapshot)},
	})

	return &SnapshotResponse{Snapshot: snapshot, UpdatedAt: time.Now()}, nil
}

// PullSnapshot returns the encrypted snapshot. If a P2P device id is supplied it
// must still be a peer — a removed device gets 403, not stale data.
func (s *Service) PullSnapshot(ctx context.Context, workspaceID, vaultID, deviceID string) (*SnapshotResponse, error) {
	v, err := s.getScoped(ctx, workspaceID, vaultID)
	if err != nil {
		return nil, err
	}
	if deviceID != "" && !hasPeer(v, deviceID) {
		return nil, apperror.New(403, "access revoked")
	}
	if len(v.Snapshot) == 0 {
		return nil, apperror.New(404, "no snapshot yet")
	}
	return &SnapshotResponse{Snapshot: v.Snapshot, UpdatedAt: v.UpdatedAt}, nil
}

// RemovePeer deauthorizes a peer (route restricts this to owner/admin).
func (s *Service) RemovePeer(ctx context.Context, workspaceID, vaultID, deviceID string) error {
	if _, err := s.getScoped(ctx, workspaceID, vaultID); err != nil {
		return err
	}
	removed, err := s.Store.RemovePeer(ctx, vaultID, deviceID)
	if err != nil {
		return apperror.ErrInternal
	}
	if !removed {
		return apperror.New(404, "peer not found")
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.peer.removed",
		TargetType:  "vault",
		TargetID:    vaultID,
		Details: 	 map[string]any{"device_id": deviceID},
	})

	return nil
}


// CreateToken adds a join token to a vault and returns its public view.
func (s *Service) CreateToken(ctx context.Context, workspaceID, vaultID string, req CreateTokenRequest) (*TokenResponse, error) {
	if _, err := s.getScoped(ctx, workspaceID, vaultID); err != nil {
		return nil, err
	}
	tok := Token{
		ID:         id.Generate("lv_join_", 12), // prefix preserved for the CLI
		Name:       req.Name,
		WrappedDEK: req.WrappedDEK,
		Verifier:   req.Verifier,
		CreatedAt:  time.Now(),
		ExpiresAt:  req.ExpiresAt, // nil = never expires
		Revoked:    false,
	}
	if err := s.Store.AddToken(ctx, vaultID, tok); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.token.created",
		TargetType:  "vault",
		TargetID:    vaultID,
		TargetName:  tok.Name,
	})

	return &TokenResponse{ID: tok.ID, Name: tok.Name, CreatedAt: tok.CreatedAt, ExpiresAt: tok.ExpiresAt}, nil
}

// ListTokens returns only active (non-revoked, non-expired) tokens, public fields only.
func (s *Service) ListTokens(ctx context.Context, workspaceID, vaultID string) (*ListTokensResponse, error) {
	v, err := s.getScoped(ctx, workspaceID, vaultID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	out := make([]TokenResponse, 0, len(v.Tokens))
	for _, t := range v.Tokens {
		if !tokenActive(t, now) {
			continue
		}
		out = append(out, TokenResponse{ID: t.ID, Name: t.Name, CreatedAt: t.CreatedAt, ExpiresAt: t.ExpiresAt})
	}
	return &ListTokensResponse{Tokens: out}, nil
}

// RevokeToken marks a token revoked so it can no longer be used to join.
func (s *Service) RevokeToken(ctx context.Context, workspaceID, vaultID, tokenID string) error {
	if _, err := s.getScoped(ctx, workspaceID, vaultID); err != nil {
		return err
	}
	ok, err := s.Store.RevokeToken(ctx, vaultID, tokenID)
	if err != nil {
		return apperror.ErrInternal
	}
	if !ok {
		return apperror.New(404, "token not found")
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: workspaceID,
		Action:      "vault.token.revoked",
		TargetType:  "vault",
		TargetID:    vaultID,
	})

	return nil
}


// Join adds the caller as a peer using a join token. Public — no JWT.
func (s *Service) Join(ctx context.Context, req JoinRequest) (*JoinResponse, error) {
	v, err := s.Store.FindVaultByTokenID(ctx, req.Token)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if v == nil {
		return nil, apperror.New(404, "invalid or expired token")
	}

	// Locate the token inside the embedded array and validate it.
	var tok *Token
	for i := range v.Tokens {
		if v.Tokens[i].ID == req.Token {
			tok = &v.Tokens[i]
			break
		}
	}
	if tok == nil || !tokenActive(*tok, time.Now()) {
		return nil, apperror.New(404, "invalid or expired token")
	}
	if !verifierMatches(tok.Verifier, req.Verifier) {
		return nil, apperror.New(404, "invalid or expired token") // wrong secret looks the same as a missing token
	}

	// Idempotent: re-joining just returns the current state.
	if hasPeer(v, req.DeviceID) {
		return &JoinResponse{
			VaultID:    v.ID,
			Snapshot:   v.Snapshot,
			Peers:      v.Peers,
			WrappedDEK: tok.WrappedDEK,
			Message:    "already a peer",
		}, nil
	}

	peer := Peer{
		DeviceID:        req.DeviceID,
		DeviceName:      req.DeviceName,
		PublicKey:       req.PublicKey,
		X25519PublicKey: req.X25519PublicKey,
		JoinedAt:        time.Now(),
	}
	if err := s.Store.AddPeer(ctx, v.ID, peer); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{
		WorkspaceID: v.WorkspaceID,
		Action:      "vault.peer.joined",
		TargetType:  "vault",
		TargetID:    v.ID,
		TargetName:  v.Name,
		Details:     map[string]any{"device_id": req.DeviceID, "device_name": req.DeviceName},
	})

	return &JoinResponse{
		VaultID:    v.ID,
		Snapshot:   v.Snapshot,
		Peers:      append(v.Peers, peer), // include the just-added peer in the reply
		WrappedDEK: tok.WrappedDEK,
	}, nil
}

// SendMessage queues one offline message with a 48h TTL.
func (s *Service) SendMessage(ctx context.Context, req SendMessageRequest) (*SendMessageResponse, error) {
	now := time.Now()
	m := &PendingMessage{
		ID:            id.Generate("msg_", 12),
		ForDeviceID:   req.ForDeviceID,
		FromDeviceID:  req.FromDeviceID,
		FromPublicKey: req.FromPublicKey,
		Payload:       req.Payload,
		CreatedAt:     now,
		ExpiresAt:     now.Add(MessageTTL),
	}
	if err := s.Store.CreateMessage(ctx, m); err != nil {
		return nil, apperror.ErrInternal
	}
	return &SendMessageResponse{ID: m.ID, Success: true}, nil
}

// GetMessages drains (returns + deletes) a device's queued messages.
func (s *Service) GetMessages(ctx context.Context, deviceID string) (*MessagesResponse, error) {
	msgs, err := s.Store.DrainMessages(ctx, deviceID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return &MessagesResponse{Messages: msgs, Count: len(msgs)}, nil
}

