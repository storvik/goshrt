{ pkgs }:
final: prev: {
  goshrt = pkgs.callPackage ./goshrt.nix { };
  goshrtc = pkgs.callPackage ./goshrtc.nix { };
}
