<div align="center">
    <a href="https://github.com/storvik/goshrt" rel="noopener">
        <img width=300px height=300px src="https://github.com/storvik/goshrt/blob/master/assets/goshrt.png" alt="goshrt_logo" />
    </a>
    <a href="https://github.com/storvik/goshrt" rel="noopener">
        <h3 align="center">goshrt</h3>
    </a>
    <p>Self hosted URL shortener<br />❗ Work in progress!</p>
</div>

<div align="center">
    <a href="https://github.com/storvik/goshrt/blob/master/LICENSE"><img src="https://img.shields.io/github/license/storvik/goshrt"></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/releases"><img src="https://img.shields.io/github/v/release/storvik/goshrt?include_prereleases"></a>
</div>

<div align="center">
    <a href="https://github.com/storvik/goshrt/actions/workflows/build.yml"><img src="https://github.com/storvik/goshrt/actions/workflows/build.yml/badge.svg" /></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/actions/workflows/gotest.yml"><img src="https://github.com/storvik/goshrt/actions/workflows/gotest.yml/badge.svg" /></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/actions/workflows/nix.yml"><img src="https://github.com/storvik/goshrt/actions/workflows/nix.yml/badge.svg" /></a>
</div>

<div align="center">
    <a href="https://goreportcard.com/report/github.com/storvik/goshrt"><img src="https://goreportcard.com/badge/github.com/storvik/goshrt" /></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/actions/workflows/lint.yml"><img src="https://github.com/storvik/goshrt/actions/workflows/lint.yml/badge.svg?branch=master" /></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/actions/workflows/testcoverage.yml"><img src="https://raw.githubusercontent.com/storvik/goshrt/badges/.badges/master/coverage.svg" /></a> &nbsp;
    <a href="https://github.com/storvik/goshrt/actions/workflows/vuln.yml"><img src="https://github.com/storvik/goshrt/actions/workflows/vuln.yml/badge.svg" /></a>
</div>

---

This is my attempt at creating a self hosted URL shortener written in Go.
The goal is to support multiple domains, cache, a simple API for creating new entries and a command line client.
Even though I use this in production bugs should be expected.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [goshrt - Self hosted URL shortener](#goshrt---self-hosted-url-shortener)
    - [Install](#install)
        - [Nix](#nix)
            - [Overlay](#overlay)
            - [Module](#module)
    - [Development](#development)
        - [Postgres](#postgres)
        - [Unit testing](#unit-testing)

<!-- markdown-toc end -->


## Install

Goshrt can easily be deployed to NixOS server using module available in the flake.nix.
Should provide instructions for traditional server and Docker.

### Nix

There are several ways to install goshrt with Nix.
Two recommended approaches is using an overlay, which will work on both NixOS and Nix with another OS, or the included NixOS module.

#### Overlay

The overlay can be added to another flake like this:

``` nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    goshrt = {
      url = "github:golang/goshrt";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ self, nixpkgs, goshrt}:
    let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        overlays = [ goshrt.overlays.default ];
      };
    in {

      # goshrt should now be available in pkgs
      devShell."x86_64-linux" = pkgs.mkShell {
        buildInputs = [
          pkgs.goshrt
          pkgs.goshrtc
        ];
      };

    };
}
```

#### Module

This is how the provided NixOS module can be used in another flake:

1. Add goshrt as input in flake.nix
2. Import goshrt service module
3. Configure service, see `module.nix` for options

The following example enables goshrt module with postgres and nginx reverse proxy.

``` nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    goshrt = {
      url = "github:golang/goshrt";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ self, nixpkgs, goshrt }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {

      nixosConfigurations.mycomputer = pkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          goshrt.nixosModules.goshrt
          ({ pkgs, ... }: {
            networking.firewall.allowedTCPPorts = [ 8080 ];

            # Needed because of `nginx.virtualHosts..enableACME`
            security.acme.acceptTerms = true;
            security.acme.defaults.email = "goshrt@example.com";

            services.goshrt = {
              enable = true;
              httpPort = 8080;
              key = "qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$";
              database = {
                enable = true; # Let goshrt module setup postgresql service
                host = "localhost";
                port = 5432;
                user = "goshrt";
                password = "trhsog";
              };
              nginx = {
                enable = true; # Enable automatic nginx proxy with SSL and LetsEncrypt certificates
                hostnames = [ "examplename1.com" "examplename2.com" ];
                extraConfig = { forceSSL = true; enableACME = true; };
              };
            };
          })
        ];
      };
    };
}
```

Password for postgres user goshrt must be set manually.
This can be achieved with:

``` shell
$ sudo -u postgres psql goshrt
> ALTER USER goshrt WITH PASSWORD 'trhsog';
```

> This includes potsgres and nginx proxy. If `config.services.goshrt.database.enable` and `config.services.goshrt.nginx.enable` is false both postgres and proxy must be setup manually.

## Development

### Postgres

When doing local development/testing postgres has to be running.
While there are several ways to achieve this, VM / docker / podman, I myself use Nix.
Spinning up a development database is very simple in Nix shell.
After installing Nix, devshell is entered through the command `nix develop`.
The following helpers for dealing with postgres is avilable:

``` shell
$ pgnix-init     # initiate database and start it
$ pgnix-start    # start database
$ pgnix-status   # check if database is running
$ pgnix-restart  # restart database
$ pgnix-stop     # stop postgresql database
$ pgnix-purge    # stop database and delete it
$ pgnix-pgcli    # start pgcli and connect to database
$ pgnix-psql     # start psql and connect to database
```

### Unit testing

Nix shell and `pgnix-` wrappers makes running unit test in a clean environment very simple.
Inside `nix develop` the following oneliner runs all unit tests:

``` shell
$ pgnix-purge && pgnix-init && go clean -testcache && go test -v ./...
```

> `go clean -testcache` ensures that all tests are run.
> Without it tests will be cached and for instance database migrations will not be run.
