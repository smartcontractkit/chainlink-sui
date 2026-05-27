{
  stdenv,
  pkgs,
}:
stdenv.mkDerivation rec {
  name = "sui-cli-${version}";
  version = "1.72.2"; # Update as needed. Should be a mainnet release version from https://github.com/MystenLabs/sui/releases

  src = if stdenv.hostPlatform.isDarwin then
    pkgs.fetchzip {
      url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-macos-arm64.tgz"; # Assume is a M1 Mac
      sha256 = "sha256-20O1Lp9qJdT+bnbgqLsQa+fe8TbrqGuOwXxVocLv3sY=";  # Should be replaced when bumping versions
      stripRoot = false;
    }
    else if stdenv.isLinux then
      pkgs.fetchzip {
        url = "https://github.com/MystenLabs/sui/releases/download/mainnet-v${version}/sui-mainnet-v${version}-ubuntu-x86_64.tgz";
        sha256 = "sha256-iOhVRR4tZJiOmQJmvcxQxEgq0cfdqtI/NYObZ+NU4GY=";  # Should be replaced when bumping versions
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
