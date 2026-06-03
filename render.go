package grafanadashboard

import (
	"encoding/json"

	errs "github.com/pleme-io/errors-go"
)

// The j* types mirror the subset of the Grafana dashboard JSON schema this
// package emits. They are private so the rendered shape is owned here; struct
// fields marshal in declaration order, which (with no Go maps in the path) makes
// RenderJSON byte-stable.

type jDashboard struct {
	UID        string   `json:"uid"`
	Title      string   `json:"title"`
	Tags       []string `json:"tags"`
	Schema     int      `json:"schemaVersion"`
	Editable   bool     `json:"editable"`
	Style      string   `json:"style"`
	Panels     []jPanel `json:"panels"`
	Templating jTempl   `json:"templating"`
	Time       jTime    `json:"time"`
	Timezone   string   `json:"timezone"`
}

type jTempl struct {
	List []struct{} `json:"list"`
}

type jTime struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type jPanel struct {
	ID         int          `json:"id"`
	Title      string       `json:"title"`
	Type       string       `json:"type"`
	GridPos    jGridPos     `json:"gridPos"`
	Datasource string       `json:"datasource,omitempty"`
	Targets    []jTarget    `json:"targets"`
	FieldCfg   jFieldConfig `json:"fieldConfig"`
}

type jGridPos struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}

type jTarget struct {
	Expr         string `json:"expr"`
	LegendFormat string `json:"legendFormat,omitempty"`
	Datasource   string `json:"datasource,omitempty"`
	RefID        string `json:"refId"`
}

type jFieldConfig struct {
	Defaults jFieldDefaults `json:"defaults"`
}

type jFieldDefaults struct {
	Unit       string       `json:"unit,omitempty"`
	Color      jColor       `json:"color"`
	Thresholds *jThresholds `json:"thresholds,omitempty"`
}

type jColor struct {
	Mode       string `json:"mode"`
	FixedColor string `json:"fixedColor,omitempty"`
}

type jThresholds struct {
	Mode  string       `json:"mode"`
	Steps []jThreshold `json:"steps"`
}

type jThreshold struct {
	Color string   `json:"color"`
	Value *float64 `json:"value"`
}

// RenderJSON renders the dashboard to Grafana dashboard JSON, themed by t. It is
// a pure function of (Dashboard, Theme): the same input yields byte-identical
// JSON every call. Panels are laid out left-to-right, wrapping at the 24-column
// grid boundary, in declaration order.
func (d *Dashboard) RenderJSON(t Theme) ([]byte, error) {
	out := jDashboard{
		UID:        d.UID,
		Title:      d.Title,
		Tags:       append([]string{"borealis", "theme:" + t.Name}, d.Tags...),
		Schema:     39,
		Editable:   true,
		Style:      "dark",
		Panels:     make([]jPanel, 0, len(d.Panels)),
		Templating: jTempl{List: []struct{}{}},
		Time:       jTime{From: "now-6h", To: "now"},
		Timezone:   "browser",
	}

	x, y := 0, 0
	rowH := 0
	for i, p := range d.Panels {
		if x+p.Width > 24 {
			x = 0
			y += rowH
			rowH = 0
		}
		jp := jPanel{
			ID:         i + 1,
			Title:      p.Title,
			Type:       string(p.Type),
			GridPos:    jGridPos{H: p.Height, W: p.Width, X: x, Y: y},
			Datasource: p.firstDatasource(d.Datasource),
			Targets:    make([]jTarget, 0, len(p.Targets)),
			FieldCfg: jFieldConfig{
				Defaults: jFieldDefaults{
					Unit:  p.Unit,
					Color: jColor{Mode: "palette-classic"},
				},
			},
		}
		for ti, tg := range p.Targets {
			jp.Targets = append(jp.Targets, jTarget{
				Expr:         tg.Expr,
				LegendFormat: tg.Legend,
				Datasource:   tg.Datasource,
				RefID:        refID(ti),
			})
		}
		if len(p.Thresholds) > 0 {
			steps := make([]jThreshold, 0, len(p.Thresholds))
			for si, th := range p.Thresholds {
				st := jThreshold{Color: t.Color(th.Role)}
				// The first step in Grafana's "absolute" mode uses a null base.
				if si != 0 {
					v := th.Value
					st.Value = &v
				}
				steps = append(steps, st)
			}
			jp.FieldCfg.Defaults.Thresholds = &jThresholds{Mode: "absolute", Steps: steps}
			// A thresholds panel colours by threshold, not the classic palette.
			jp.FieldCfg.Defaults.Color = jColor{Mode: "thresholds"}
		}
		out.Panels = append(out.Panels, jp)
		x += p.Width
		if p.Height > rowH {
			rowH = p.Height
		}
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return nil, errs.Wrap(err, "grafanadashboard: marshal dashboard JSON", errs.WithCode("dashboard_render"))
	}
	b = append(b, '\n')
	return b, nil
}

// firstDatasource returns the panel's effective datasource: the first target's
// override if set, else the dashboard default.
func (p Panel) firstDatasource(dflt string) string {
	for _, tg := range p.Targets {
		if tg.Datasource != "" {
			return tg.Datasource
		}
	}
	return dflt
}

// refID maps a target index to Grafana's A/B/C… refId convention.
func refID(i int) string {
	if i < 26 {
		return string(rune('A' + i))
	}
	return "A" + itoa(i)
}
