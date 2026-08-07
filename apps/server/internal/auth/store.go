package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store handles all MongoDB operations for auth
type Store struct {
	users    *mongo.Collection
	sessions *mongo.Collection
}

// NewStore creates store and sets up MongoDB indexes for fast queries
func NewStore(db *mongo.Database) *Store {
	s := &Store{
		users:    db.Collection("users"),
		sessions: db.Collection("sessions"),
	}

	ctx := context.TODO()

	// Unique index on email — prevents two accounts with same email
	s.users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},   // bson.D = ordered document, 1 = ascending
		Options: options.Index().SetUnique(true),
	})

	// TTL index — MongoDB auto-deletes documents when expires_at passes (no cleanup cron needed)
	s.sessions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	s.sessions.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}}, // speeds up "find all sessions for user"
	})

	return s
}

// --- User operations ---
func (s *Store) CreateUser(ctx context.Context, user *User) error {
	_, err := s.users.InsertOne(ctx, user)
	return err
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	// bson.M = unordered map, like {"email": "test@example.com"}
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil // not found
	}
	return &user, err // &user = pointer to user, so caller gets the actual data
}

func (s *Store) FindUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := s.users.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

func (s *Store) FindUserByOAuth(ctx context.Context, provider, oauthID string) (*User, error) {
	var user User
	err := s.users.FindOne(ctx, bson.M{
		"oauth_provider": provider,
		"oauth_id":       oauthID,
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}

// UpdateUser updates specific fields — $set changes only listed fields, leaves rest untouched
func (s *Store) UpdateUser(ctx context.Context, id string, fields bson.M) error {
	_, err := s.users.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": fields})
	return err
}

// ------------------------ Session operations ---
func (s *Store) CreateSession(ctx context.Context, session *Session) error {
	_, err := s.sessions.InsertOne(ctx, session)
	return err
}

// FindSessionByTokenHash finds active session — $gt means "greater than" (not expired)
func (s *Store) FindSessionByTokenHash(ctx context.Context, tokenHash string) (*Session, error) {
	var session Session
	err := s.sessions.FindOne(ctx, bson.M{
		"refresh_token_hash": tokenHash,
		"expires_at":         bson.M{"$gt": time.Now()},
	}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &session, err
}

func (s *Store) DeleteSession(ctx context.Context, id string) error {
	_, err := s.sessions.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *Store) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := s.sessions.DeleteOne(ctx, bson.M{"refresh_token_hash": tokenHash})
	return err
}

// FindSessionsByUserID returns all active sessions for user — used by "manage sessions" page
func (s *Store) FindSessionsByUserID(ctx context.Context, userID string) ([]Session, error) {
	cursor, err := s.sessions.Find(ctx, bson.M{
		"user_id":    userID,
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	var sessions []Session
	err = cursor.All(ctx, &sessions) // cursor.All reads all results into the slice
	return sessions, err
}

// DeleteSessionsByUserIDExcept deletes all sessions except one — used after password change
func (s *Store) DeleteSessionsByUserIDExcept(ctx context.Context, userID, exceptID string) error {
	// $ne = "not equal" — keeps the current session, deletes all others
	_, err := s.sessions.DeleteMany(ctx, bson.M{
		"user_id": userID,
		"_id":     bson.M{"$ne": exceptID},
	})
	return err
}

// DeleteSessionsByDeviceID removes all sessions for one device — used on device revoke.
func (s *Store) DeleteSessionsByDeviceID(ctx context.Context, deviceID string) error {
	_, err := s.sessions.DeleteMany(ctx, bson.M{"device_id": deviceID})
	return err
}
