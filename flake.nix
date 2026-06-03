# flake.nix — grafana-dashboard-go (GSDS Biblioteca) via substrate's go-library-flake.
# vendorHash OMITTED → spec-sourced (__from-spec__); clean nix build lands once
# errors-go is published. Pre-publish proof is `go test` (green).
{
  description = "grafana-dashboard-go — typed Grafana dashboard model → dashboard JSON, borealis-themed";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    substrate = {
      # Published repo uses: url = "github:pleme-io/substrate";
      url = "git+file:///Users/drzzln/code/github/pleme-io/substrate?ref=feat/go-pattern-parity";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs @ { self, nixpkgs, substrate, ... }:
    (import substrate.goLibraryFlakeBuilder { inherit nixpkgs; }) {
      name = "grafana-dashboard-go";
      version = "0.1.0";
      src = self;
      repo = "pleme-io/grafana-dashboard-go";
    };
}
