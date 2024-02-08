{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    # (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.go_1_21
    pkgs.gomod2nix
    pkgs.govulncheck
    pkgs.golangci-lint
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
