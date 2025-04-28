{
  description = "Goshrt flake file";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ self, nixpkgs }:
    let
      linuxSystems = [ "x86_64-linux" "aarch64-linux" ];
      darwinSystems = [ "aarch64-darwin" "x86_64-darwin" ];
      forAllSystems = function:
        nixpkgs.lib.genAttrs (linuxSystems ++ darwinSystems)
          (system:
            function (import nixpkgs {
              inherit system;
              overlays = [ ];
              config = { };
            }));
      version = "0.3.0";
    in
    {

      overlays = forAllSystems (pkgs: {
        default = import ./nix/overlay.nix { inherit pkgs; };
      });

      nixosModules.goshrt = import ./nix/module.nix self;

      # Small container which is meant to test nixos module in CI
      nixosConfigurations.container = nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          self.nixosModules.goshrt
          ./nix/module_test.nix
        ];
      };

      packages = forAllSystems (pkgs: {
        goshrt = pkgs.callPackage ./nix/goshrt.nix { inherit version; };
        goshrtc = pkgs.callPackage ./nix/goshrtc.nix { inherit version; };
        default = pkgs.callPackage ./nix/goshrt.nix { inherit version; };
      });

      devShells = forAllSystems (pkgs: {
        default = import ./shell.nix { inherit pkgs; };
      });
    };
}
