{
  description = "Goshrt flake file";

  inputs = {

    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

  };

  outputs = inputs@{ self, nixpkgs, flake-utils, gitignore, gomod2nix }:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        overlays = [
          gomod2nix.overlays.default
          self.overlays.default
        ];
      };
      inherit (gitignore.lib) gitignoreSource;
    in
    {

      overlays.default = import ./overlay.nix { inherit pkgs; };

      nixosModules.goshrt = import ./module.nix self;

      packages."x86_64-linux" = {
        goshrt = pkgs.callPackage ./goshrt.nix { };
        goshrtc = pkgs.callPackage ./goshrtc.nix { };

        default = pkgs.callPackage ./goshrt.nix { };
      };

      devShells."x86_64-linux".default = import ./shell.nix { inherit pkgs; };
    };
}
