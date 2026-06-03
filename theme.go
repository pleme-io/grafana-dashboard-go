package grafanadashboard

// Role is a semantic colour intent, mirroring the borealis design system's Role
// vocabulary. A consumer maps a domain state (a threshold band, a panel
// emphasis) onto a Role; the [Theme] decides the actual hex, so every fleet
// dashboard reads the same and colours are never hand-authored (the ishou /
// borealis "consume tokens, never hex" rule).
type Role int

const (
	// RoleNeutral is inactive / baseline (Polar Night surface tone).
	RoleNeutral Role = iota
	// RolePrimary is the headline / primary-metric colour (Frost).
	RolePrimary
	// RoleSuccess is a healthy band (Aurora green).
	RoleSuccess
	// RoleInfo is informational (Frost cyan).
	RoleInfo
	// RoleWarning is a degraded band (Aurora yellow).
	RoleWarning
	// RoleDanger is a breached band (Aurora red).
	RoleDanger
)

// Theme is the resolved token bundle: a Role → hex map plus the Grafana base
// colour names this theme prefers. It is pure data so a dashboard renders the
// same regardless of where the theme came from.
type Theme struct {
	// Name identifies the theme in the rendered dashboard metadata.
	Name string
	// roles maps each semantic Role to its hex colour.
	roles map[Role]string
}

// Color resolves a [Role] to its hex string under this theme. An unknown Role
// falls back to the neutral surface tone so rendering never panics.
func (t Theme) Color(r Role) string {
	if hex, ok := t.roles[r]; ok {
		return hex
	}
	return t.roles[RoleNeutral]
}

// Tundra is the pleme-io fleet-default Grafana theme: the Nord Aurora family
// bound to semantic roles, matching the borealis Tundra theme. These are the
// canonical Nord hex tokens (nordtheme.com) — the same palette ishou binds.
func Tundra() Theme {
	return Theme{
		Name: "tundra",
		roles: map[Role]string{
			RoleNeutral: "#4C566A", // nord3  — Polar Night surface
			RolePrimary: "#88C0D0", // nord8  — Frost primary
			RoleInfo:    "#8FBCBB", // nord7  — Frost cyan
			RoleSuccess: "#A3BE8C", // nord14 — Aurora green
			RoleWarning: "#EBCB8B", // nord13 — Aurora yellow
			RoleDanger:  "#BF616A", // nord11 — Aurora red
		},
	}
}

// Nord is an alias of [Tundra] under the upstream palette name — same tokens,
// provided so a consumer can name the source palette explicitly.
func Nord() Theme {
	t := Tundra()
	t.Name = "nord"
	return t
}
