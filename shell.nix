{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.gomod2nix
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
