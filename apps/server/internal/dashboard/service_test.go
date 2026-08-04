package dashboard

import (
	"testing"
	"time"
)

func TestParseRange(t *testing.T) {
	tests := []struct {
		raw      string
		wantKey  string
		wantDays int
		wantErr  bool
	}{
		{raw: "", wantKey: "1y", wantDays: 365},
		{raw: "1y", wantKey: "1y", wantDays: 365},
		{raw: "365d", wantKey: "1y", wantDays: 365},
		{raw: "7d", wantKey: "7d", wantDays: 7},
		{raw: "30d", wantKey: "30d", wantDays: 30},
		{raw: "14d", wantErr: true},
		{raw: "week", wantErr: true},
	}
	for _, tt := range tests {
		key, days, err := parseRange(tt.raw)
		if tt.wantErr {
			if err == nil {
				t.Fatalf("parseRange(%q): expected error", tt.raw)
			}
			continue
		}
		if err != nil {
			t.Fatalf("parseRange(%q): unexpected error: %v", tt.raw, err)
		}
		if key != tt.wantKey || days != tt.wantDays {
			t.Fatalf("parseRange(%q) = (%q, %d); want (%q, %d)", tt.raw, key, days, tt.wantKey, tt.wantDays)
		}
	}
}

func TestFillSeriesEmptyDays(t *testing.T) {
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	out := fillSeries(from, 7, nil)
	if len(out) != 7 {
		t.Fatalf("len = %d; want 7", len(out))
	}
	if out[0].Date != "2026-08-01" || out[6].Date != "2026-08-07" {
		t.Fatalf("date range = %s..%s; want 2026-08-01..2026-08-07", out[0].Date, out[6].Date)
	}
	for _, p := range out {
		if p.Total != 0 {
			t.Fatalf("day %s total = %d; want 0", p.Date, p.Total)
		}
		if p.ByPrefix == nil {
			t.Fatalf("day %s by_prefix is nil", p.Date)
		}
		if len(p.ByPrefix) != 0 {
			t.Fatalf("day %s by_prefix = %v; want empty", p.Date, p.ByPrefix)
		}
	}
}

func TestFillSeriesMergesPrefixes(t *testing.T) {
	from := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	rows := []seriesRow{
		rowAt(from, "vault", 3),
		rowAt(from, "member", 2),
		rowAt(from.AddDate(0, 0, 2), "workspace", 1),
	}
	out := fillSeries(from, 3, rows)
	if len(out) != 3 {
		t.Fatalf("len = %d; want 3", len(out))
	}
	if out[0].Total != 5 || out[0].ByPrefix["vault"] != 3 || out[0].ByPrefix["member"] != 2 {
		t.Fatalf("day0 = %+v; want total 5 vault=3 member=2", out[0])
	}
	if out[1].Total != 0 || len(out[1].ByPrefix) != 0 {
		t.Fatalf("day1 = %+v; want empty", out[1])
	}
	if out[2].Total != 1 || out[2].ByPrefix["workspace"] != 1 {
		t.Fatalf("day2 = %+v; want total 1 workspace=1", out[2])
	}
}

func rowAt(day time.Time, prefix string, count int64) seriesRow {
	var r seriesRow
	r.ID.Date = day
	r.ID.Prefix = prefix
	r.Count = count
	return r
}
