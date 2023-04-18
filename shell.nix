{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.gomod2nix
  ] ++ (import ./pgnix.nix { inherit pkgs; });
}
