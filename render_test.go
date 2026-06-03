package grafanadashboard

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	errs "github.com/pleme-io/errors-go"
)

func sampleDashboard(t *testing.T) *Dashboard {
	t.Helper()
	d, err := New("service-red",
		WithTitle("Service RED"),
		WithTag("service"),
		WithDatasource("prometheus"),
		WithPanel(Panel{
			Title:   "Request rate",
			Type:    TimeSeries,
			Unit:    "reqps",
			Targets: []Target{{Expr: "rate(http_requests_total[5m])", Legend: "{{method}}"}},
		}),
		WithPanel(Panel{
			Title:   "Error ratio",
			Type:    Gauge,
			Unit:    "percentunit",
			Targets: []Target{{Expr: "rate(http_requests_total{code=~\"5..\"}[5m])"}},
			Thresholds: []Threshold{
				{Value: 0, Role: RoleSuccess},
				{Value: 0.01, Role: RoleWarning},
				{Value: 0.05, Role: RoleDanger},
			},
		}))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return d
}

func TestNew_Validation(t *testing.T) {
	cases := []struct {
		name     string
		uid      string
		opts     []Option
		wantCode string
	}{
		{"missing uid", "", nil, "dashboard_invalid"},
		{"panel no title", "d", []Option{WithPanel(Panel{Targets: []Target{{Expr: "up"}}})}, "dashboard_invalid"},
		{"panel no targets", "d", []Option{WithPanel(Panel{Title: "x"})}, "dashboard_invalid"},
		{"valid empty", "d", nil, ""},
		{"valid with panel", "d", []Option{WithPanel(Panel{Title: "x", Targets: []Target{{Expr: "up"}}})}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.uid, tc.opts...)
			if got := errs.CodeOf(err); got != tc.wantCode {
				t.Fatalf("code = %q, want %q (err=%v)", got, tc.wantCode, err)
			}
		})
	}
}

func TestNew_AppliesDefaults(t *testing.T) {
	d, err := New("d", WithPanel(Panel{Title: "p", Targets: []Target{{Expr: "up"}}}))
	if err != nil {
		t.Fatal(err)
	}
	if d.Title != "d" {
		t.Errorf("title default = %q, want uid", d.Title)
	}
	p := d.Panels[0]
	if p.Type != TimeSeries || p.Width != 12 || p.Height != 8 {
		t.Errorf("panel defaults not applied: %+v", p)
	}
}

func TestRenderJSON_ValidAndShaped(t *testing.T) {
	d := sampleDashboard(t)
	raw, err := d.RenderJSON(Tundra())
	if err != nil {
		t.Fatalf("RenderJSON: %v", err)
	}
	// Must be valid JSON.
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("rendered dashboard is not valid JSON: %v", err)
	}
	cases := []struct {
		name  string
		check func() bool
	}{
		{"uid", func() bool { return doc["uid"] == "service-red" }},
		{"title", func() bool { return doc["title"] == "Service RED" }},
		{"has panels", func() bool { ps, _ := doc["panels"].([]any); return len(ps) == 2 }},
		{"theme tag", func() bool { return bytes.Contains(raw, []byte("theme:tundra")) }},
		{"borealis tag", func() bool { return bytes.Contains(raw, []byte("borealis")) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.check() {
				t.Errorf("check %q failed", tc.name)
			}
		})
	}
}

func TestRenderJSON_ThresholdsThemed(t *testing.T) {
	d := sampleDashboard(t)
	th := Tundra()
	raw, _ := d.RenderJSON(th)
	out := string(raw)
	cases := []struct {
		role Role
		hex  string
	}{
		{RoleSuccess, th.Color(RoleSuccess)},
		{RoleWarning, th.Color(RoleWarning)},
		{RoleDanger, th.Color(RoleDanger)},
	}
	for _, tc := range cases {
		t.Run(tc.hex, func(t *testing.T) {
			if !strings.Contains(out, tc.hex) {
				t.Errorf("themed threshold colour %q (role %d) missing from output", tc.hex, tc.role)
			}
		})
	}
	// The first threshold step carries a null base value (Grafana convention).
	if !strings.Contains(out, "\"value\": null") {
		t.Error("first threshold step must have a null base value")
	}
}

// TestByteStable is the load-bearing invariant: deterministic JSON every render.
func TestByteStable(t *testing.T) {
	d := sampleDashboard(t)
	for _, th := range []Theme{Tundra(), Nord()} {
		t.Run(th.Name, func(t *testing.T) {
			first, _ := d.RenderJSON(th)
			for i := 0; i < 16; i++ {
				next, _ := d.RenderJSON(th)
				if !bytes.Equal(first, next) {
					t.Fatalf("RenderJSON(%s) not byte-stable", th.Name)
				}
			}
		})
	}
}

func TestRenderJSON_GridLayoutWraps(t *testing.T) {
	// Three panels of width 12 → row1: [0..12),[12..24); row2: [0..12).
	d, err := New("grid",
		WithPanel(Panel{Title: "a", Width: 12, Height: 8, Targets: []Target{{Expr: "1"}}}),
		WithPanel(Panel{Title: "b", Width: 12, Height: 8, Targets: []Target{{Expr: "2"}}}),
		WithPanel(Panel{Title: "c", Width: 12, Height: 8, Targets: []Target{{Expr: "3"}}}))
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := d.RenderJSON(Tundra())
	var doc struct {
		Panels []struct {
			Title   string             `json:"title"`
			GridPos struct{ X, Y int } `json:"gridPos"`
		} `json:"panels"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	want := map[string][2]int{"a": {0, 0}, "b": {12, 0}, "c": {0, 8}}
	for _, p := range doc.Panels {
		w := want[p.Title]
		if p.GridPos.X != w[0] || p.GridPos.Y != w[1] {
			t.Errorf("panel %q at (%d,%d), want (%d,%d)", p.Title, p.GridPos.X, p.GridPos.Y, w[0], w[1])
		}
	}
}

func TestRefID(t *testing.T) {
	cases := map[int]string{0: "A", 1: "B", 25: "Z", 26: "A26"}
	for i, want := range cases {
		if got := refID(i); got != want {
			t.Errorf("refID(%d) = %q, want %q", i, got, want)
		}
	}
}
