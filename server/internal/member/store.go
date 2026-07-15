package member

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store handles MongoDB for the member domain. It owns "workspace_invites" and
// shares "memberships"; it also reads "users" and "workspaces" for display/email.
type Store struct {
	client      *mongo.Client     // needed to start sessions for the join transaction
	memberships *mongo.Collection // shared with the workspace domain
	invites     *mongo.Collection // owned here
	users       *mongo.Collection // read-only: display fields + email lookup
	workspaces  *mongo.Collection // read-only: workspace name for the invite email
}


// NewStore creates the store and sets up indexes for workspace_invites.
// (memberships/users/workspaces indexes are created by their owning domains.)
func NewStore(db *mongo.Database) *Store {
	s := &Store{
		client:      db.Client(), // reuse the already-connected client for sessions
		memberships: db.Collection("memberships"),
		invites:     db.Collection("workspace_invites"),
		users:       db.Collection("users"),
		workspaces:  db.Collection("workspaces"),
	}

	ctx := context.TODO()

	// TTL — MongoDB auto-deletes an invite once expires_at passes (no cleanup cron).
	s.invites.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}}, // 1 = ascending
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	// token_hash — join looks an invite up by its hashed token.
	s.invites.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "token_hash", Value: 1}},
	})
	// (workspace_id, status) — list pending invites for a workspace quickly.
	s.invites.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "workspace_id", Value: 1}, {Key: "status", Value: 1}},
	})

	return s
}


// GetMembership returns the caller/target's membership, or nil if they are not a member.
func (s *Store) GetMembership(ctx context.Context, workspaceID, userID string) (*Membership, error) {
	var m Membership
	err := s.memberships.FindOne(ctx, bson.M{
		"workspace_id": workspaceID,
		"user_id":      userID,
	}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &m, err
}

// RoleOf returns the caller's role in a workspace, or "" if not a member.
// This makes *Store satisfy middleware.MembershipChecker (reused RequireRole from Step 3).
func (s *Store) RoleOf(ctx context.Context, workspaceID, userID string) (string, error) {
	m, err := s.GetMembership(ctx, workspaceID, userID)
	if err != nil {
		return "", err
	}
	if m == nil {
		return "", nil
	}
	return m.Role, nil
}

// MembershipExists is a cheap "is this user already a member?" check.
func (s *Store) MembershipExists(ctx context.Context, workspaceID, userID string) (bool, error) {
	err := s.memberships.FindOne(ctx, bson.M{
		"workspace_id": workspaceID,
		"user_id":      userID,
	}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// ListMemberships returns every membership row for a workspace (for the members list).
func (s *Store) ListMemberships(ctx context.Context, workspaceID string) ([]Membership, error) {
	cursor, err := s.memberships.Find(ctx, bson.M{"workspace_id": workspaceID})
	if err != nil {
		return nil, err
	}
	var members []Membership
	err = cursor.All(ctx, &members) // read all matching docs into the slice
	return members, err
}

// UpdateMemberRole changes one member's role ($set leaves other fields untouched).
func (s *Store) UpdateMemberRole(ctx context.Context, workspaceID, userID, role string) error {
	_, err := s.memberships.UpdateOne(ctx,
		bson.M{"workspace_id": workspaceID, "user_id": userID},
		bson.M{"$set": bson.M{"role": role}},
	)
	return err
}

// DeleteMembership removes one member from a workspace.
func (s *Store) DeleteMembership(ctx context.Context, workspaceID, userID string) error {
	_, err := s.memberships.DeleteOne(ctx, bson.M{
		"workspace_id": workspaceID,
		"user_id":      userID,
	})
	return err
}

// PendingInviteExists reports whether this email already has a live pending invite here.
func (s *Store) PendingInviteExists(ctx context.Context, workspaceID, email string) (bool, error) {
	err := s.invites.FindOne(ctx, bson.M{
		"workspace_id": workspaceID,
		"email":        email,
		"status":       StatusPending,
		"expires_at":   bson.M{"$gt": time.Now()}, // ignore ones the TTL hasn't swept yet
	}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CreateInvite inserts one pending invite.
func (s *Store) CreateInvite(ctx context.Context, inv *Invite) error {
	_, err := s.invites.InsertOne(ctx, inv)
	return err
}

// ListPendingInvites returns all live pending invites for a workspace.
func (s *Store) ListPendingInvites(ctx context.Context, workspaceID string) ([]Invite, error) {
	cursor, err := s.invites.Find(ctx, bson.M{
		"workspace_id": workspaceID,
		"status":       StatusPending,
		"expires_at":   bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	var invites []Invite
	err = cursor.All(ctx, &invites)
	return invites, err
}

// FindInviteByTokenHash looks up an invite by its hashed token, or nil if none.
func (s *Store) FindInviteByTokenHash(ctx context.Context, hash string) (*Invite, error) {
	var inv Invite
	err := s.invites.FindOne(ctx, bson.M{"token_hash": hash}).Decode(&inv)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &inv, err
}

// FindInviteByID looks up an invite by id, or nil if none (used to cancel).
func (s *Store) FindInviteByID(ctx context.Context, id string) (*Invite, error) {
	var inv Invite
	err := s.invites.FindOne(ctx, bson.M{"_id": id}).Decode(&inv)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &inv, err
}

// DeleteInvite removes one invite (cancel).
func (s *Store) DeleteInvite(ctx context.Context, id string) error {
	_, err := s.invites.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AcceptInvite atomically inserts the new membership and flips the invite to accepted.
// Both commit or neither does — same StartSession/WithTransaction pattern as workspace.
func (s *Store) AcceptInvite(ctx context.Context, inv *Invite, mem *Membership) error {
	session, err := s.client.StartSession() // a session is required to run a transaction
	if err != nil {
		return err
	}
	defer session.EndSession(ctx) // always release the session

	// The ctx passed into fn carries the session — pass it to every op so they join the txn.
	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		if _, err := s.memberships.InsertOne(ctx, mem); err != nil {
			return nil, err
		}
		_, err := s.invites.UpdateOne(ctx,
			bson.M{"_id": inv.ID},
			bson.M{"$set": bson.M{"status": StatusAccepted}},
		)
		if err != nil {
			return nil, err // membership insert is rolled back automatically
		}
		return nil, nil // returning nil error commits the transaction
	})
	return err
}

// FindUserByEmail returns the user with this email, or nil if no account exists yet.
// Read-only view of the auth domain's "users" collection (display + already-member check).
func (s *Store) FindUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	var u UserInfo
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &u, err
}

// FindUsersByIDs returns display info for a batch of user ids ($in = one query for the list).
func (s *Store) FindUsersByIDs(ctx context.Context, ids []string) ([]UserInfo, error) {
	cursor, err := s.users.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var users []UserInfo
	err = cursor.All(ctx, &users)
	return users, err
}

// FindWorkspaceName returns just the workspace's display name (for the invite email).
func (s *Store) FindWorkspaceName(ctx context.Context, workspaceID string) (string, error) {
	// projection: pull back only the "name" field, nothing else.
	var doc struct {
		Name string `bson:"name"`
	}
	err := s.workspaces.FindOne(ctx, bson.M{"_id": workspaceID},
		options.FindOne().SetProjection(bson.M{"name": 1}),
	).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return "", nil
	}
	return doc.Name, err
}