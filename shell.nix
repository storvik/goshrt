{ pkgs }:

pkgs.mkShell
{
  buildInputs = [
    pkgs.go
    pkgs.govulncheck
    pkgs.golangci-lint
  ] ++ (import ./nix/pgnix.nix { inherit pkgs; });
}
