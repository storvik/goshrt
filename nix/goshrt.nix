{ stdenv
, callPackage
, go
, lib
, buildGoModule
, version
}:

buildGoModule {
  pname = "goshrt";
  inherit version;
  src = ./..;
  # modRoot = ./..;
  subPackages = [ "cmd/goshrt" ];
  vendorHash = "sha256-63ube2xpfqkfCbiAO79BgIEH6JgVkmnAg4HUURsZjLI=";
}
