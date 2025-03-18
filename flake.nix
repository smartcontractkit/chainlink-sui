{
  description = "Sui integration";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = inputs @ {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {
        inherit system;
        # A predicate to allow specific packages: redpanda-rpk has an unfree license (‘bsl11’), refusing to evaluate.
        config.allowUnfreePredicate = pkg:
          builtins.elem (nixpkgs.lib.getName pkg) [
            "redpanda-rpk"
          ];
      };
    in rec {
      devShell = pkgs.callPackage ./shell.nix {};
    });
}
