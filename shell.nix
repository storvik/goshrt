{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    # disabled as it does not pickup go version from go.mod
    # (pkgs.mkGoEnv { pwd = ./.; })
    pkgs.go_1_21
    pkgs.gomod2nix
    pkgs.govulncheck
    pkgs.golangci-lint
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
