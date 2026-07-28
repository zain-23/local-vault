package workspace

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store handles all mongoDB operations for workspace domain
type Store struct {
	client			*mongo.Client
	workspaces		*mongo.Collection
	memberships 	*mongo.Collection
}

// NewStore creates the store and sets up indexes
func NewStore(db *mongo.Database) *Store {
	s := &Store{
		client: 		db.Client(),  // reuse the already connected client for sessions
		workspaces: 	db.Collection("workspaces"),
		memberships: 	db.Collection("memberships"),
	}

	ctx := context.TODO()

	// unique slug - no two workspaces share a URL handle
	s.workspaces.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: 		bson.D{{Key: "slug", Value: 1}}, // 1 = ascending
		Options: 	options.Index().SetUnique(true),
	})

	// owner_id - speed up "workspaces i own"
	s.workspaces.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: 		bson.D{{Key: "owner_id", Value: 1}},
	})

	// unique (workspace_id, user_id) - a user can join a workspace at most once
	s.memberships.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: 		bson.D{{Key: "workspace_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	s.memberships.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})

	return s
}


// SlugExists reports whether a workspace already uses this slug
func (s *Store) SlugExists(ctx context.Context, slug string) (bool, error) {
	err := s.workspaces.FindOne(ctx, bson.M{"slug": slug}).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil // free to use
	}
	if err != nil {
		return false, err
	}
	return true, nil // taken
}


// CreateWorkspaceWithOwner inserts the workspace + owner membership automatically in one transaction
func (s *Store) CreateWorkspaceWithOwner(ctx context.Context, ws *Workspace, mem *Membership) error {
	session, err := s.client.StartSession() // session is required for run a transaction

	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		if _, err := s.workspaces.InsertOne(ctx, ws); err != nil {
			return nil, err
		}

		if _, err := s.memberships.InsertOne(ctx, mem); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

// GetMembership returns the caller's membership in a workspace, or nil if they are not a member
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
// This method makes *Store satisfy middleware.MembershipChecker (Task 6).
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

// FindWorkspaceByID returns one workspace or nil if it does not exist
func (s *Store) FindWorkspaceByID(ctx context.Context, id string) (*Workspace, error) {
	var w Workspace
	err := s.workspaces.FindOne(ctx, bson.M{"_id": id}).Decode(&w)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &w, err
}

// FindMembershipsByUserID returns every membership row for a user (to list their workspaces)
func (s *Store) FindMembershipsByUserID(ctx context.Context, userID string) ([]Membership, error) {
	cursor, err := s.memberships.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var members []Membership
	err = cursor.All(ctx, &members) // read all matching docs into the slice
	return members, err
}

// FindWorkspacesByIDs returns all workspaces whose _id is in the list — one query for the list endpoint
func (s *Store) FindWorkspacesByIDs(ctx context.Context, ids []string) ([]Workspace, error) {
	// $in matches any document whose _id appears in ids
	cursor, err := s.workspaces.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var ws []Workspace
	err = cursor.All(ctx, &ws)
	return ws, err
}

// UpdateWorkspace sets the given fields on a workspace ($set leaves other fields untouched)
func (s *Store) UpdateWorkspace(ctx context.Context, id string, fields bson.M) error {
	_, err := s.workspaces.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

// DeleteWorkspaceCascade removes the workspace and all its memberships atomically
func (s *Store) DeleteWorkspaceCascade(ctx context.Context, id string) error {
	session, err := s.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(ctx context.Context) (any, error) {
		if _, err := s.workspaces.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
			return nil, err
		}
		// DeleteMany removes every membership pointing at this workspace
		if _, err := s.memberships.DeleteMany(ctx, bson.M{"workspace_id": id}); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}