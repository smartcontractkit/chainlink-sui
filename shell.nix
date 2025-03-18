{
  stdenv,
  pkgs,
  lib,
}:
# juno requires building with clang, not gcc
(pkgs.mkShell.override {stdenv = pkgs.clangStdenv;}) {
  buildInputs = with pkgs;
    [
      # Go 1.23 + tools
      go_1_23
      gopls

      # Keep adding as needed

      # Sui CLI custom derivation
      (pkgs.callPackage ./sui.nix {})
    ];
}
