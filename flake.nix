{
  description = "Goshrt flake file";

  inputs = {

    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    nix2container.url = "github:nlewo/nix2container";

    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

  };

  outputs = inputs@{ self, nixpkgs, gitignore, nix2container }:
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
      nix2containerPkgs = nix2container.packages.x86_64-linux;
      inherit (gitignore.lib) gitignoreSource;
    in
    rec {

      # packages."x86_64-linux".goshrt-dev-container = nix2containerPkgs.nix2container.buildImage {
      #   name = "goshrt-dev-container";
      #   config = {
      #     Cmd = [ "postgres" ];
      #     WorkingDir = "/data";
      #   };
      # };

      devShell."x86_64-linux" =
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
            pkgs.glibcLocales
            pkgs.postgresql
            pkgs.pgcli
          ];
        };

    };
}
