package device

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store handles MongoDB for the device domain. It owns both collections outright.
type Store struct {
	requests *mongo.Collection // device_auth_requests — in-flight logins
	devices  *mongo.Collection // devices — authorized CLI installations
}

// NewStore creates the store and sets up indexes for both collections.
func NewStore(db *mongo.Database) *Store {
	s := &Store{
		requests: db.Collection("device_auth_requests"),
		devices:  db.Collection("devices"),
	}

	ctx := context.TODO()

	// TTL — MongoDB auto-deletes a request once expires_at passes (no cleanup cron).
	s.requests.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}}, // 1 = ascending
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	// user_code — the browser looks a request up by it. Unique: two live requests
	// sharing a code would let the wrong CLI be approved.
	s.requests.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "user_code", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	// device_code_hash — the poll endpoint's lookup key. Unique for the same reason.
	s.requests.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "device_code_hash", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	// (user_id, workspace_id) — list a user's devices in a workspace quickly.
	s.devices.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "workspace_id", Value: 1}},
	})

	return s
}


// ---------------- Auth request operations ----------------

// CreateAuthRequest inserts a new pending login attempt.
func (s *Store) CreateAuthRequest(ctx context.Context, req *AuthRequest) error {
	_, err := s.requests.InsertOne(ctx, req)
	return err
}

// findAuthRequest is the shared lookup — filter differs, decoding doesn't.
// Returns (nil, nil) when nothing matches, so callers distinguish "absent" from "broken".
func (s *Store) findAuthRequest(ctx context.Context, filter bson.M) (*AuthRequest, error) {
	var req AuthRequest
	err := s.requests.FindOne(ctx, filter).Decode(&req)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil // not found is not an error — the caller decides the status code
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *Store) FindAuthRequestByUserCode(ctx context.Context, userCode string) (*AuthRequest, error) {
	return s.findAuthRequest(ctx, bson.M{"user_code": userCode})
}

// FindAuthRequestByDeviceCodeHash looks a request up by the hashed secret (CLI poll side).
func (s *Store) FindAuthRequestByDeviceCodeHash(ctx context.Context, hash string) (*AuthRequest, error) {
	return s.findAuthRequest(ctx, bson.M{"device_code_hash": hash})
}

// ApproveAuthRequest stamps the approver, workspace and new device onto the request.
// The status filter makes this a no-op if the request already left "pending" —
// two browser tabs clicking Approve must not create two devices.
func (s *Store) ApproveAuthRequest(ctx context.Context, id, userID, workspaceID, deviceID string) error {
	res, err := s.requests.UpdateOne(ctx,
		bson.M{"_id": id, "status": StatusPending},
		bson.M{"$set": bson.M{
			"status":       StatusApproved,
			"user_id":      userID,
			"workspace_id": workspaceID,
			"device_id":    deviceID,
		}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("request was not pending") // already approved, denied, or gone
	}
	return nil
}

// DenyAuthRequest marks the request denied so the CLI stops polling.
func (s *Store) DenyAuthRequest(ctx context.Context, id string) error {
	_, err := s.requests.UpdateOne(ctx,
		bson.M{"_id": id, "status": StatusPending},
		bson.M{"$set": bson.M{"status": StatusDenied}},
	)
	return err
}

// ConsumeAuthRequest flips consumed=false → true and reports whether THIS call won.
// The consumed:false filter is what makes token collection single-use: two racing
// polls both read status=approved, but only one UpdateOne matches a document.
func (s *Store) ConsumeAuthRequest(ctx context.Context, id string) (bool, error) {
	res, err := s.requests.UpdateOne(ctx,
		bson.M{"_id": id, "consumed": false},
		bson.M{"$set": bson.M{"consumed": true}},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount == 1, nil // false = someone already collected the tokens
}

// ---------------- Device operations ----------------

// CreateDevice inserts an authorized device.
func (s *Store) CreateDevice(ctx context.Context, d *Device) error {
	_, err := s.devices.InsertOne(ctx, d)
	return err
}

// ListDevices returns one user's devices in one workspace, newest first.
// Both filters matter: a user must not be able to enumerate a teammate's devices.
func (s *Store) ListDevices(ctx context.Context, userID, workspaceID string) ([]Device, error) {
	opts := options.Find().SetSort(bson.D{{Key: "authorized_at", Value: -1}}) // -1 = descending
	cur, err := s.devices.Find(ctx, bson.M{"user_id": userID, "workspace_id": workspaceID}, opts)
	if err != nil {
		return nil, err
	}
	// All decodes the whole cursor into the slice and closes it for us.
	var devices []Device
	if err := cur.All(ctx, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

// FindDeviceByID returns one device, or (nil, nil) if it doesn't exist.
func (s *Store) FindDeviceByID(ctx context.Context, id string) (*Device, error) {
	var d Device
	err := s.devices.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// DeleteDevice removes a device row. Its sessions are killed separately by the
// service, via auth.Service.RevokeDeviceSessions.
func (s *Store) DeleteDevice(ctx context.Context, id string) error {
	_, err := s.devices.DeleteOne(ctx, bson.M{"_id": id})
	return err
}