package audit

import (
	"context"
	"log"
	"time"

	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/id"
	"github.com/zain-23/local-vault/server/internal/common/pagination"
	"github.com/zain-23/local-vault/server/internal/common/reqctx"
)

// Recorder is the one method other domains depend on. Declaring it here and
// injecting it lets each domain stay decoupled and lets tests pass a no-op.
// Record returns NOTHING: auditing must never fail a user's action.
type Recorder interface {
	Record(ctx context.Context, e Entry)
}

// Service implements Recorder and the read side.
type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

// Record writes one event. It reads actor/ip/device from the request context
// (set by middleware.Auth) and swallows any error — a failed audit write is
// logged, never propagated.
func (s *Service) Record(ctx context.Context, e Entry) {
	info := reqctx.From(ctx) // zero-value on public routes — that's fine

	ev := &Event{
		ID:          id.Generate("evt_", 12),
		WorkspaceID: e.WorkspaceID,
		ActorID:     info.ActorID,
		Action:      e.Action,
		TargetType:  e.TargetType,
		TargetID:    e.TargetID,
		TargetName:  e.TargetName,
		Details:     e.Details,
		DeviceID:    info.DeviceID,
		IP:          info.IP,
		CreatedAt:   time.Now(),
	}

	// Detach from the request lifecycle: the handler may return before this insert
	// finishes, and a cancelled request ctx must not drop the audit record.
	writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.store.Insert(writeCtx, ev); err != nil {
		log.Printf("⚠️ audit: failed to record %q: %v", e.Action, err)
	}
}

// List returns one filtered, paginated page of a workspace's events.
func (s *Service) List(ctx context.Context, workspaceID string, q ListQuery) (*pagination.Page[EventResponse], error) {
	f, err := parseFilter(q)
	if err != nil {
		return nil, err
	}
	items, total, err := s.store.ListPaginated(ctx, workspaceID, f, q.Params)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return &pagination.Page[EventResponse]{
		Items: items,
		Meta:  pagination.NewMeta(q.Params, total),
	}, nil
}

// Export returns every matching event for CSV download.
func (s *Service) Export(ctx context.Context, workspaceID string, q ListQuery) ([]EventResponse, error) {
	f, err := parseFilter(q)
	if err != nil {
		return nil, err
	}
	events, err := s.store.ExportAll(ctx, workspaceID, f)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return events, nil
}

// parseFilter validates the query and parses the RFC3339 date bounds.
func parseFilter(q ListQuery) (Filter, error) {
	f := Filter{
		Action:       q.Action,
		ActionPrefix: q.ActionPrefix,
		ActorID:      q.ActorID,
		DeviceID:     q.DeviceID,
	}
	if q.From != "" {
		t, err := time.Parse(time.RFC3339, q.From)
		if err != nil {
			return Filter{}, apperror.New(400, "from must be an RFC3339 timestamp")
		}
		f.From = &t
	}
	if q.To != "" {
		t, err := time.Parse(time.RFC3339, q.To)
		if err != nil {
			return Filter{}, apperror.New(400, "to must be an RFC3339 timestamp")
		}
		f.To = &t
	}
	return f, nil
}
