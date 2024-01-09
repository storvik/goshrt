{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.gomod2nix
    pkgs.govulncheck
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
