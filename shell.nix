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
    echo "Setting up clean Go environment (disabling GVM)..."
    # Unset GVM environment leakage
    unset GOROOT
    unset GOPATH
    unset GOTOOLDIR
    # Add Nix-provided Go binary path to ensure consistency
    export PATH=$(go env GOROOT)/bin:$PATH
    # Optional: lock Go toolchain version if needed
    export GOTOOLCHAIN=go1.24.2+auto
    # Debug info
    echo "Using Go at: $(which go)"
    go version
    bun --version
    # use upstream golangci-lint config from core Chainlink repository, overriding the local prefixes
    alias golint="golangci-lint run --config <(curl -sSL https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/.golangci.yml | yq e '.formatters.settings.goimports.local-prefixes = [\"github.com/smartcontractkit/chainlink-ton\"]' -) --path-mode \"abs\""
    echo ""
    echo "You can lint your code with:"
    echo "    cd relayer && golint ./..."
    echo "    cd integration-tests && golint ./..."
    echo ""
  '';
})
