package grafanadashboard

import (
	errs "github.com/pleme-io/errors-go"
)

// PanelType is the Grafana panel renderer kind.
type PanelType string

const (
	// TimeSeries is the line/area time-series panel.
	TimeSeries PanelType = "timeseries"
	// Stat is the single-number stat panel.
	Stat PanelType = "stat"
	// Gauge is the gauge panel.
	Gauge PanelType = "gauge"
	// Table is the table panel.
	Table PanelType = "table"
)

// Dashboard is the one typed declaration of a Grafana dashboard. Construct with
// [New] and functional options; the constructor applies defaults and validates
// invariants in one place.
type Dashboard struct {
	// UID is the stable dashboard UID. Required (seeds the JSON `uid`).
	UID string
	// Title is the human dashboard title. Defaults to the UID.
	Title string
	// Tags are dashboard tags, kept in declaration order.
	Tags []string
	// Datasource is the default datasource name applied to targets that do not
	// set their own. Empty leaves Grafana's default.
	Datasource string
	// Panels are the dashboard's panels, laid out in declaration order.
	Panels []Panel
}

// Panel is one dashboard panel.
type Panel struct {
	// Title is the panel heading.
	Title string
	// Type selects the panel renderer. Defaults to [TimeSeries].
	Type PanelType
	// Unit is the value unit (Grafana unit id, e.g. "reqps", "s", "percent").
	Unit string
	// Width is the panel grid width in 24-column units (default 12).
	Width int
	// Height is the panel grid height in row units (default 8).
	Height int
	// Targets are the panel's queries.
	Targets []Target
	// Thresholds are the value bands; their colours resolve from the theme Role.
	Thresholds []Threshold
}

// Target is one panel query.
type Target struct {
	// Expr is the query expression (e.g. a PromQL string).
	Expr string
	// Legend is the legend format. Optional.
	Legend string
	// Datasource overrides the dashboard default for this target. Optional.
	Datasource string
}

// Threshold is one value band. The colour is not a hex string but a semantic
// [Role] resolved from the [Theme] at render time — so dashboards never carry
// hand-authored colours.
type Threshold struct {
	// Value is the band's lower bound. The first threshold conventionally has
	// Value 0 and is the base band.
	Value float64
	// Role is the semantic colour intent for this band.
	Role Role
}

// Option configures a [Dashboard] during [New].
type Option func(*Dashboard)

// WithTitle sets the dashboard title.
func WithTitle(title string) Option { return func(d *Dashboard) { d.Title = title } }

// WithTag appends a dashboard tag.
func WithTag(tag string) Option { return func(d *Dashboard) { d.Tags = append(d.Tags, tag) } }

// WithDatasource sets the default datasource.
func WithDatasource(ds string) Option { return func(d *Dashboard) { d.Datasource = ds } }

// WithPanel appends a panel.
func WithPanel(p Panel) Option { return func(d *Dashboard) { d.Panels = append(d.Panels, p) } }

// New builds and validates a [Dashboard]. It returns a typed, code-carrying
// error (code "dashboard_invalid") when an invariant is violated.
func New(uid string, opts ...Option) (*Dashboard, error) {
	d := &Dashboard{UID: uid}
	for _, o := range opts {
		o(d)
	}
	if d.UID == "" {
		return nil, errs.New("grafanadashboard: uid is required", errs.WithCode("dashboard_invalid"))
	}
	if d.Title == "" {
		d.Title = d.UID
	}
	for i := range d.Panels {
		p := &d.Panels[i]
		if p.Title == "" {
			return nil, errs.New("grafanadashboard: panel "+itoa(i)+" has no title", errs.WithCode("dashboard_invalid"))
		}
		if p.Type == "" {
			p.Type = TimeSeries
		}
		if p.Width == 0 {
			p.Width = 12
		}
		if p.Height == 0 {
			p.Height = 8
		}
		if len(p.Targets) == 0 {
			return nil, errs.New("grafanadashboard: panel "+p.Title+" has no targets", errs.WithCode("dashboard_invalid"))
		}
	}
	return d, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
