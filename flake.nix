{
  description = "Goshrt flake file";

  inputs = {

    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    nix2container.url = "github:nlewo/nix2container";

    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:tweag/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

  };

  outputs = inputs@{ self, nixpkgs, flake-utils, gitignore, nix2container, gomod2nix }:
    let
      systems = [ "x86_64-linux" "aarch64" ];
    in
    flake-utils.lib.eachSystem systems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
        nix2containerPkgs = nix2container.packages.x86_64-linux;
        inherit (gitignore.lib) gitignoreSource;
        pversion = "0.0.1";
      in
      rec {

        packages = {
          goshrt = pkgs.buildGoApplication {
            pname = "goshrt";
            version = pversion;
            pwd = ./.;
            src = ./.;
            subPackages = [ "cmd/goshrt" ];
            modules = ./gomod2nix.toml;
            ldflags = ''
              -w -s -X 'github.com/storvik/goshrt/version.GitVersion=${pversion}'
                    -X 'github.com/storvik/goshrt/version.GitCommit=${self.shortRev or "dirty"}'
            '';
            doCheck = false;
          };

          goshrtc = pkgs.buildGoApplication {
            pname = "goshrtc";
            version = pversion;
            pwd = ./.;
            src = ./.;
            subPackages = [ "cmd/goshrtc" ];
            modules = ./gomod2nix.toml;
            ldflags = ''
              -w -s -X 'github.com/storvik/goshrt/version.GitVersion=${pversion}'
                    -X 'github.com/storvik/goshrt/version.GitCommit=${self.shortRev or "dirty"}'
            '';
            doCheck = false;
          };

          defaultPackage = packages.goshrt;
        };

        devShell =
          let
            tmp = ".devshell";

            db = "goshrt";
            user = "goshrt";
            passwd = "trhsog";
            port = "6000";
            pgdata = ".devshell/db";
          in
          pkgs.mkShell {
            PGDATA = pgdata;

            buildInputs = [
              (pkgs.mkGoEnv { pwd = ./.; })
              pkgs.glibcLocales
              pkgs.postgresql
              pkgs.pgcli
              pkgs.gomod2nix
              (pkgs.writeScriptBin "pgnix-init" ''
                initdb -D ${pgdata} -U postgres
                pg_ctl -D ${pgdata} -l ${pgdata}/postgres.log  -o "-p ${port} -k /tmp -i" start
                createdb --port=${port} --host=localhost --username=postgres -O postgres ${db}
                psql -d postgres -U postgres -h localhost -p ${port} -c "create user ${user} with encrypted password '${passwd}';"
                psql -d postgres -U postgres -h localhost -p ${port} -c "grant all privileges on database ${db} to ${user};"
              '')
              (pkgs.writeScriptBin "pgnix-start" ''
                pg_ctl -D ${pgdata} -l ${pgdata}/postgres.log  -o "-p ${port} -k /tmp -i" start
              '')
              (pkgs.writeScriptBin "pgnix-pgcli" ''
                PGPASSWORD=${passwd} pgcli -h localhost -p 6000 -U goshrt
              '')
              (pkgs.writeScriptBin "pgnix-psql" ''
                PGPASSWORD=${passwd} psql -d ${db} -U ${user} -h localhost -p ${port}
              '')
              (pkgs.writeScriptBin "pgnix-status" ''
                pg_ctl -D ${pgdata} status
              '')
              (pkgs.writeScriptBin "pgnix-restart" ''
                pg_ctl -D ${pgdata} restart
              '')
              (pkgs.writeScriptBin "pgnix-stop" ''
                pg_ctl -D ${pgdata} stop
              '')
              (pkgs.writeScriptBin "pgnix-purge" ''
                pg_ctl -D ${pgdata} stop
                rm -rf .devshell/db
              '')
            ];
          };
      });
}
