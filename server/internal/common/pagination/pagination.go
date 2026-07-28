package pagination

import "math"

// Defaults and bounds for the page/limit query params.
const (
	DefaultPage  = 1   // page is 1-based
	DefaultLimit = 20  // rows per page when limit is omitted
	MaxLimit     = 100 // hard cap so a client can't ask for the whole table
)

// Params are the pagination query fields. Embed this in a request struct so
// Fiber's QueryParser fills page/limit and the validator checks the bounds.
// omitempty = a missing param decodes to 0 and skips validation (defaults apply).
type Params struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}

// Normalize replaces zero-values (an omitted param) with the defaults. Call this
// after validation, before using the params in a query.
func (p *Params) Normalize() {
	if p.Page == 0 {
		p.Page = DefaultPage
	}
	if p.Limit == 0 {
		p.Limit = DefaultLimit
	}
}

// Skip is the Mongo $skip offset for the current page (0-based).
func (p Params) Skip() int { return (p.Page - 1) * p.Limit }

// Meta is the pagination block returned to the client alongside the items.
type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`       // total rows after filtering
	TotalPages int   `json:"total_pages"` // ceil(total/limit); 0 when total is 0
}

// Page is the generic paginated response body: the rows plus their meta.
type Page[T any] struct {
	Items []T  `json:"items"`
	Meta  Meta `json:"meta"`
}

// NewMeta builds the meta block from the (normalized) params and a total count.
func NewMeta(p Params, total int64) Meta {
	pages := 0
	if total > 0 && p.Limit > 0 {
		pages = int(math.Ceil(float64(total) / float64(p.Limit)))
	}
	return Meta{Page: p.Page, Limit: p.Limit, Total: total, TotalPages: pages}
}
