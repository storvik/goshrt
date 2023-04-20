{ stdenv
, callPackage
, go
, lib
, buildGoApplication
}:

buildGoApplication rec {
  pname = "goshrt";
  version = "0.0.1";
  pwd = ./..;
  src = ./..;
  subPackages = [ "cmd/goshrt" ];
  modules = ./../gomod2nix.toml;
  # ldflags = "-w -s -X 'github.com/storvik/goshrt/version.GitVersion=${version}'";
  doCheck = false;

  meta = {
    description = "Self hosted URL shortener server written in Go";
    homepage = "https://github.com/storvik/goshrt";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.petterstorvik ];
  };
}
