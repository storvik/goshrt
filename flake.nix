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
      version = "0.3.0";
      inherit (gitignore.lib) gitignoreSource;
    in
    {

      overlays.default = import ./nix/overlay.nix { inherit pkgs; };

      nixosModules.goshrt = import ./nix/module.nix self;

      # Small container which is meant to test nixos module in CI
      nixosConfigurations.container = nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          self.nixosModules.goshrt
          ./nix/module_test.nix
        ];
      };

      packages."x86_64-linux" = {
        goshrt = pkgs.callPackage ./nix/goshrt.nix { inherit version; };
        goshrtc = pkgs.callPackage ./nix/goshrtc.nix { inherit version; };
        default = pkgs.callPackage ./nix/goshrt.nix { inherit version; };
      };

      devShells."x86_64-linux".default = import ./shell.nix { inherit pkgs; };
    };
}
