# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-06-03

### Added
- Initial: a typed `Dashboard` of `Panel`s (`Target`s + `Threshold`s) rendered
  byte-stable to Grafana dashboard JSON via `RenderJSON(Theme)`. Borealis-aligned
  `Tundra()`/`Nord()` themes binding semantic `Role`s to the Nord Aurora palette,
  auto grid layout on the 24-column grid, and typed code-carrying errors via
  `errors-go` (`dashboard_invalid`/`dashboard_render`).
