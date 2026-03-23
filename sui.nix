{
  stdenv,
  pkgs,
}:
stdenv.mkDerivation rec {
  name = "sui-cli-${version}";
  version = "1.65.2"; # Update as needed. Should be a mainnet release version from https://github.com/MystenLabs/sui/releases

  src = if stdenv.hostPlatform.isDarwin then
    pkgs.fetchzip {
      url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-macos-arm64.tgz"; # Assume is a M1 Mac
      sha256 = "sha256-5wwht1qMp68k3bHIDIyixPFteJ7cW3oVKk/5GwcICjM=";  # Should be replaced when bumping versions
      stripRoot = false;
    }
    else if stdenv.isLinux then
      pkgs.fetchzip {
        url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-ubuntu-x86_64.tgz";
        sha256 = "sha256-eK8nvQfV2W8oQzOr1kNVc/3TABwvDG+6SCm/MIQpz5I=";  # Should be replaced when bumping versions
        stripRoot = false;
      }
    else
      builtins.throw "Unsupported system";

  sourceRoot = ".";

  # No build needed since we're using a prebuilt binary archive.
  buildPhase = "true";
  installPhase = ''
    mkdir -p $out/bin
    mv source/sui $out/bin/
  '';
}
