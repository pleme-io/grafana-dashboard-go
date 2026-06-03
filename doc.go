// Package grafanadashboard is the fleet's typed Grafana dashboard model — a
// dashboard authored as a typed [Dashboard] of [Panel]s (with [Target]s and
// [Threshold]s) and rendered, borealis-themed, to Grafana dashboard JSON. So
// dashboards are types, not hand-edited JSON blobs that drift.
//
// # Why
//
// Grafana dashboards live as 1000-line JSON files that nobody hand-edits
// safely: a copy-pasted panel keeps a stale datasource, thresholds disagree
// across panels, the colour palette is whatever the last editor clicked. The
// cure is to author the dashboard as a typed value and render the JSON
// mechanically with a consistent (borealis/Tundra Aurora) theme.
//
// # Shape (Law 1, §3.5)
//
//	d, err := grafanadashboard.New("service-red",
//	    grafanadashboard.WithTitle("Service RED"),
//	    grafanadashboard.WithPanel(grafanadashboard.Panel{
//	        Title:   "Request rate",
//	        Type:    grafanadashboard.TimeSeries,
//	        Targets: []grafanadashboard.Target{{Expr: "rate(http_requests_total[5m])"}},
//	        Thresholds: []grafanadashboard.Threshold{
//	            {Value: 0, Role: grafanadashboard.RoleSuccess},
//	            {Value: 100, Role: grafanadashboard.RoleWarning},
//	        },
//	    }))
//	jsonBytes, err := d.RenderJSON(grafanadashboard.Tundra())
//
// RenderJSON is a pure function of the [Dashboard] + [Theme]: the same input
// yields byte-identical JSON every time (deterministic field order, panel grid
// laid out in declaration order, threshold colours resolved from the theme).
//
// # Theme (borealis-aligned)
//
// [Tundra] and [Nord] return the same Aurora-family palette the borealis design
// system binds to semantic roles. Thresholds and panel colours reference
// semantic [Role]s, not raw hex, so every fleet dashboard reads the same.
//
// # Weight (Law 6)
//
// Pure standard library (encoding/json) plus errors-go for typed, code-carrying
// errors. The theme palette is baked in as typed tokens; it does not import the
// heavy borealis charm stack — a dashboard renderer needs the hex tokens, not a
// terminal renderer.
package grafanadashboard
