# grafana-dashboard-go

A typed Grafana dashboard model (panels / targets / thresholds) rendered,
borealis-themed, to Grafana dashboard JSON — so dashboards are types, not
hand-edited 1000-line JSON blobs that drift.

## What

A typed [`Dashboard`](./dashboard.go) of [`Panel`](./dashboard.go)s (each with
[`Target`](./dashboard.go)s and [`Threshold`](./dashboard.go)s) rendered by one
pure verb:

```go
jsonBytes, err := d.RenderJSON(grafanadashboard.Tundra())
```

`RenderJSON` is a pure function of the `Dashboard` and the `Theme`: the same
input yields **byte-identical JSON** every call (deterministic field order,
panels auto-laid-out on the 24-column grid in declaration order, threshold
colours resolved from the theme).

## Why

Grafana dashboards live as JSON nobody hand-edits safely: copy-pasted panels
keep stale datasources, thresholds disagree across panels, the palette is
whatever the last editor clicked. Authoring the dashboard as a typed value and
rendering the JSON mechanically — with a single fleet-consistent theme — removes
the whole class of drift.

## Theme (borealis-aligned)

Thresholds and panel colours reference semantic `Role`s
(`RoleSuccess`/`RoleWarning`/`RoleDanger`/…), never raw hex. `Tundra()` and
`Nord()` return the Nord Aurora palette the borealis design system binds to
those roles, so every fleet dashboard reads the same.

## Install

```
go get github.com/pleme-io/grafana-dashboard-go
```

## Usage

```go
d, err := grafanadashboard.New("service-red",
    grafanadashboard.WithTitle("Service RED"),
    grafanadashboard.WithDatasource("prometheus"),
    grafanadashboard.WithPanel(grafanadashboard.Panel{
        Title:   "Request rate",
        Type:    grafanadashboard.TimeSeries,
        Unit:    "reqps",
        Targets: []grafanadashboard.Target{{Expr: "rate(http_requests_total[5m])"}},
    }),
    grafanadashboard.WithPanel(grafanadashboard.Panel{
        Title: "Error ratio",
        Type:  grafanadashboard.Gauge,
        Targets: []grafanadashboard.Target{{Expr: "..."}},
        Thresholds: []grafanadashboard.Threshold{
            {Value: 0, Role: grafanadashboard.RoleSuccess},
            {Value: 0.01, Role: grafanadashboard.RoleWarning},
            {Value: 0.05, Role: grafanadashboard.RoleDanger},
        },
    }))
if err != nil { return errs.Exit(err) } // typed code-carrying error (dashboard_invalid)

j, _ := d.RenderJSON(grafanadashboard.Tundra())
os.WriteFile("service-red.json", j, 0o644)
```

## Configuration

None — a pure library. Callers that read the dashboard shape from config use
`shikumi-go` and pass the fields via the functional options; a `FromConfig`
bridge over a `shikumi`-loaded `Dashboard`-shaped struct is the natural
extension for config-driven dashboards.

## Release

Pull-model (Go modules): an annotated `vX.Y.Z` tag is the release; pkg.go.dev
indexes it. See the GSDS module delivery FSM.
