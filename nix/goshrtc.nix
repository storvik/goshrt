{ stdenv
, callPackage
, go
, lib
, buildGoModule
, version
}:

buildGoModule {
  pname = "goshrtc";
  inherit version;
  src = ./..;
  # modRoot = ./..;
  subPackages = [ "cmd/goshrtc" ];
  vendorHash = "sha256-63ube2xpfqkfCbiAO79BgIEH6JgVkmnAg4HUURsZjLI=";
  doCheck = false;

  meta = {
    description = "Client for goshrt - self hosted URL shortener";
    homepage = "https://github.com/storvik/goshrt";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.petterstorvik ];
  };
}
