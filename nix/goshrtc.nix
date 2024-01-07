{ stdenv
, callPackage
, go
, lib
, buildGoApplication
, version
}:

buildGoApplication rec {
  inherit version;
  pname = "goshrtc";
  pwd = ./..;
  src = ./..;
  subPackages = [ "cmd/goshrtc" ];
  modules = ./../gomod2nix.toml;
  # ldflags = "-w -s -X 'github.com/storvik/goshrt/version.GitVersion=${version}'";
  doCheck = false;

  meta = {
    description = "Client for goshrt - self hosted URL shortener";
    homepage = "https://github.com/storvik/goshrt";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.petterstorvik ];
  };
}
