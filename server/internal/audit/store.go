package audit

import (
	"context"
	"regexp"

	"github.com/zain-23/local-vault/server/internal/common/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type Store struct {
	events  	*mongo.Collection
	users  	*mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	s := &Store{
		events: db.Collection("audit_events"),
		users:  db.Collection("users"),
	}

	ctx := context.TODO()

	// (workspace_id, created_at desc) — the primary query: a workspace's log newest-first.
	s.events.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "workspace_id", Value: 1}, {Key: "created_at", Value: -1}},
	})
	// actor_id — "everything user X did".
	s.events.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "actor_id", Value: 1}},
	})

	return s
}


// Insert appends one event. Called only by the recorder.
func (s *Store) Insert(ctx context.Context, e *Event) error {
	_, err := s.events.InsertOne(ctx, e)
	return err
}

// buildMatch turns a Filter into the $match stage (shared by list + export).
func buildMatch(workspaceID string, f Filter) bson.M {
	m := bson.M{"workspace_id": workspaceID}
	if f.Action != "" {
		m["action"] = f.Action
	}
	if f.ActionPrefix != "" {
		// anchored, escaped regex: "vault" → ^vault  (QuoteMeta blocks regex injection)
		m["action"] = bson.M{"$regex": "^" + regexp.QuoteMeta(f.ActionPrefix)}
	}
	if f.ActorID != "" {
		m["actor_id"] = f.ActorID
	}
	if f.DeviceID != "" {
		m["device_id"] = f.DeviceID
	}
	// created_at range — only add the bounds that were supplied.
	if f.From != nil || f.To != nil {
		rng := bson.M{}
		if f.From != nil {
			rng["$gte"] = *f.From
		}
		if f.To != nil {
			rng["$lte"] = *f.To
		}
		m["created_at"] = rng
	}
	return m
}


// actorLookup are the shared stages that join each event to its actor's name.
// preserveNullAndEmptyArrays keeps events with no actor (public join) in the result.
func actorLookup() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "actor_id",
			"foreignField": "_id",
			"as":           "actor",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$actor", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$addFields", Value: bson.M{"actor_name": "$actor.name"}}},
		{{Key: "$project", Value: bson.M{"actor": 0}}}, // drop the joined doc, keep actor_name
	}
}

// ListPaginated returns one page of events (newest first, actor name joined) plus
// the total count after filtering — one aggregation via $facet, like member list.
func (s *Store) ListPaginated(ctx context.Context, workspaceID string, f Filter, p pagination.Params) ([]EventResponse, int64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: buildMatch(workspaceID, f)}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}, // -1 = newest first
		{{Key: "$facet", Value: bson.M{
			"items": append(mongo.Pipeline{
				{{Key: "$skip", Value: p.Skip()}},
				{{Key: "$limit", Value: p.Limit}},
			}, actorLookup()...), // join actor names only on the page, not the whole match
			"total": bson.A{bson.M{"$count": "count"}},
		}}},
	}

	cursor, err := s.events.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}

	var facet []struct {
		Items []EventResponse `bson:"items"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}
	if err := cursor.All(ctx, &facet); err != nil {
		return nil, 0, err
	}
	if len(facet) == 0 {
		return []EventResponse{}, 0, nil
	}
	items := facet[0].Items
	if items == nil {
		items = []EventResponse{} // never nil — the client expects an array
	}
	var total int64
	if len(facet[0].Total) > 0 {
		total = facet[0].Total[0].Count
	}
	return items, total, nil
}


// ExportAll returns every matching event (no pagination) for CSV export, newest first.
func (s *Store) ExportAll(ctx context.Context, workspaceID string, f Filter) ([]EventResponse, error) {
	pipeline := append(mongo.Pipeline{
		{{Key: "$match", Value: buildMatch(workspaceID, f)}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
	}, actorLookup()...)

	cursor, err := s.events.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var out []EventResponse
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
