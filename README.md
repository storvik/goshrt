# goshrt - Self hosted URL shortener

[![build](https://github.com/storvik/goshrt/actions/workflows/build.yml/badge.svg)](https://github.com/storvik/goshrt/actions/workflows/build.yml)
[![go test](https://github.com/storvik/goshrt/actions/workflows/gotest.yml/badge.svg)](https://github.com/storvik/goshrt/actions/workflows/gotest.yml)

> Work in progress!

This is my attempt at creating a self hosted URL shortener written in Go.
The goal is to support multiple domains, cache, a simple API for creating new entries and a command line client.

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
      url = "path:/home/storvik/developer/golang/goshrt";
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

``` nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    goshrt = {
      url = "path:/home/storvik/developer/golang/goshrt";
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

            services.postgresql = {
              enable = true;
              ensureDatabases = [ "goshrt" ];
              ensureUsers = [
                {
                  name = "goshrt";
                  ensurePermissions = {
                    "DATABASE goshrt" = "ALL PRIVILEGES";
                  };
                }
              ];
            };

            services.goshrt = {
              enable = true;
              httpPort = 8080;
              key = "qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$";
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

> This does not setup forwarding. Typically a nginx reverse proxy, or similar, should be used to forward requests to goshrt.

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
