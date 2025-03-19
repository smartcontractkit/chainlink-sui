{
  stdenv,
  pkgs,
  lib,
}:
(pkgs.pkgs.mkShell {
  buildInputs = with pkgs;
    [
      # Go 1.23 + tools
      go_1_23
      gopls

      # Keep adding as needed

      # Sui CLI custom derivation
      (pkgs.callPackage ./sui.nix {})
    ]
    ++ lib.optionals stdenv.hostPlatform.isDarwin [
      libiconv
    ];
})
