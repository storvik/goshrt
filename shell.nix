{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.gomod2nix
    pkgs.govulncheck
    pkgs.golangci-lint
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
