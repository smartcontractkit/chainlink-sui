{
  stdenv,
  pkgs,
  lib,
}:
(pkgs.pkgs.mkShell {
  buildInputs = with pkgs;
    [
      # Development
      go-task
      golangci-lint

      # Required for Sui CLI (Move compilation)
      git
      # Go 1.24 + tools
      go_1_24
      gopls

      # Keep adding as needed

      # Sui CLI custom derivation
      (pkgs.callPackage ./sui.nix {})
    ]
    ++ lib.optionals stdenv.hostPlatform.isDarwin [
      libiconv
    ];
})
