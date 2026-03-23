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
      # Go 1.25 + tools
      go_1_25
      gopls
      mockgen

      # bun for typescript examples
      bun

      # Keep adding as needed

      # Sui CLI custom derivation
      (pkgs.callPackage ./sui.nix {})
    ]
    ++ lib.optionals stdenv.hostPlatform.isDarwin [
      libiconv
    ];

  shellHook = ''
    >&2 echo "Setting up clean Go environment (disabling GVM)..."
    # Unset GVM environment leakage
    unset GOROOT
    unset GOPATH
    unset GOTOOLDIR
    # Use the nix-provided Go toolchain; disable Go's auto-download behavior
    export GOTOOLCHAIN=local
    # Add Nix-provided Go binary path to ensure consistency
    export PATH=$(go env GOROOT)/bin:$PATH
    # Debug info
    >&2 echo "Using Go at: $(which go)"
    >&2 go version
    >&2 bun --version
    # use upstream golangci-lint config from core Chainlink repository, overriding the local prefixes
    alias golint="golangci-lint run --config <(curl -sSL https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/.golangci.yml | yq e '.formatters.settings.goimports.local-prefixes = [\"github.com/smartcontractkit/chainlink-ton\"]' -) --path-mode \"abs\""
    >&2 echo ""
    >&2 echo "You can lint your code with:"
    >&2 echo "    cd relayer && golint ./..."
    >&2 echo "    cd integration-tests && golint ./..."
    >&2 echo ""
  '';
})
