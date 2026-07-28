package account

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store struct {
	users    *mongo.Collection
	sessions *mongo.Collection
}

// NewStore reuses the existing collections; auth already created their indexes.
func NewStore(db *mongo.Database) *Store {
	return &Store{
		users:    db.Collection("users"),
		sessions: db.Collection("sessions"),
	}
}

// FindUserByID returns the user, or (nil, nil) if missing.
func (s *Store) FindUserByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := s.users.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateUser sets the given fields ($set leaves the rest untouched).
func (s *Store) UpdateUser(ctx context.Context, id string, fields bson.M) error {
	_, err := s.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

// ConsumeBackupCode pulls one hashed code from the array. Returns whether it matched —
// $pull only modifies the doc if the code was present, so it's single-use by construction.
func (s *Store) ConsumeBackupCode(ctx context.Context, userID, hash string) (bool, error) {
	res, err := s.users.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"backup_codes": hash}},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount == 1, nil
}

// ListSessions returns the user's active sessions.
func (s *Store) ListSessions(ctx context.Context, userID string) ([]Session, error) {
	cur, err := s.sessions.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var out []Session
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FindSessionByID returns one session, or (nil, nil) if missing.
func (s *Store) FindSessionByID(ctx context.Context, id string) (*Session, error) {
	var sess Session
	err := s.sessions.FindOne(ctx, bson.M{"_id": id}).Decode(&sess)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sess, nil
}

// DeleteSession removes one session.
func (s *Store) DeleteSession(ctx context.Context, id string) error {
	_, err := s.sessions.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteSessionsExcept removes all of a user's sessions except one ($ne = not equal).
func (s *Store) DeleteSessionsExcept(ctx context.Context, userID, exceptID string) error {
	_, err := s.sessions.DeleteMany(ctx, bson.M{
		"user_id": userID,
		"_id":     bson.M{"$ne": exceptID},
	})
	return err
}
