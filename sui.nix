{
  stdenv,
  pkgs,
}:
stdenv.mkDerivation rec {
  name = "sui-cli-${version}";
  version = "1.44.3"; # Update as needed. Should be a mainnet release version

  dontUnpack = true;

  src = if stdenv.hostPlatform.isDarwin then
    pkgs.fetchurl {
      url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-macos-x86_64.tgz";
      sha256 = "sha256-bC21R4Zu2+XYK5hJb1T8wPNXWZZ/jiJw/fIOJB6QhDI=";  # Replace with the actual SHA256 for macOS asset
    }
    else if stdenv.isLinux then
      pkgs.fetchurl {
        url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-ubuntu-x86_64.tgz";
        sha256 = "sha256-bC21R4Zu2+XYK5hJb1T8wPNXWZZ/jiJw/fIOJB6QhDI=";  # Replace with the actual SHA256 for Ubuntu asset
      }
    else
      builtins.throw "Unsupported system";

  # No build needed since we're using a prebuilt binary archive.
  buildPhase = "true";
  installPhase = ''
    mkdir -p extracted
    tar -xzf $src -C extracted
    # Locate the sui binary inside the extracted contents.
    sui_bin=$(find extracted -type f -name sui -executable | head -n 1)
    if [ -z "$sui_bin" ]; then
      echo "Error: sui binary not found in the archive" >&2
      exit 1
    fi
    mkdir -p $out/bin
    cp "$sui_bin" $out/bin/
  '';
}