self:
{ config
, lib
, pkgs
, ...
}:

with lib;

let
  goshrtPkgs = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
  cfg = config.services.goshrt;
  pg = config.services.postgresql;
  goshrt = "${goshrtPkgs}/bin/goshrt --config ${serverConfig}";
  serverConfig = pkgs.writeText "goshrtconfig.toml" ''
    [server]
    key = "${cfg.key}"
    port = ":${(toString cfg.httpPort)}"
    slug_length = ${(toString cfg.slugLength)}

    [database]
    db = "${cfg.database.name}"
    user = "${cfg.database.user}"
    password = "${cfg.database.password}"
    address = "${cfg.database.host}:${(toString cfg.database.port)}"
  '';
in

{
  options.services.goshrt = {
    enable = mkOption {
      default = false;
      type = types.bool;
      description = lib.mdDoc "Enable goshrt service.";
    };
    httpPort = mkOption {
      default = 3000;
      type = types.int;
      description = lib.mdDoc "HTTP listen port.";
    };
    key = mkOption {
      type = types.str;
      # TODO: Implement keyFile as an alternative
      description = "Secret master key, will be visible in nix store.";
    };
    slugLength = mkOption {
      default = 8;
      type = types.int;
      description = lib.mdDoc "Slug length for generating random slugs.";
    };
    database = {
      enablePostgres = mkOption {
        default = false;
        type = types.bool;
        description = lib.mdDoc ''
          Enable postgres setup. If not enabled `services.postgresql`
          must be setup manually or non local postgres instance can
          be used.
        '';
      };
      host = mkOption {
        type = types.str;
        default = "127.0.0.1";
        description = lib.mdDoc "Database host address.";
      };
      port = mkOption {
        type = types.port;
        default = pg.port;
        defaultText = literalExpression ''
          config.${options.services.postgresql.port}
        '';
        description = lib.mdDoc "Database host port.";
      };
      name = mkOption {
        type = types.str;
        default = "goshrt";
        description = lib.mdDoc "Database name.";
      };
      user = mkOption {
        type = types.str;
        default = "goshrt";
        description = lib.mdDoc "Database user.";
      };

      # TODO: Implement passwordFile as an alternative
      password = mkOption {
        type = types.str;
        default = "trhsog";
        description = lib.mdDoc ''
          The password corresponding to {option}`database.user`.
          Warning: this is stored in cleartext in the Nix store!
        '';
      };
    };
  };

  config = mkIf cfg.enable {

    systemd.services.goshrt = {
      description = "goshrt - self hosted URL shortener";
      after = [ "network.target" "postgresql.service" ];
      path = [ goshrtPkgs ];
      wantedBy = [ "multi-user.target" ];
      preStart = "${goshrt} database migrate";
      serviceConfig = {
        Type = "simple";
        ExecStart = "${goshrt}";
        Restart = "always";
      };
    };

    services.postgresql = mkIf cfg.database.enablePostgres {
      enable = true;
      ensureDatabases = [ cfg.database.name ];
      ensureUsers = [
        {
          name = cfg.database.user;
          ensurePermissions = {
            "DATABASE ${cfg.database.name}" = "ALL PRIVILEGES";
          };
        }
      ];
    };

    warnings =
      optional cfg.database.enablePostgres ''
        config.services.goshrt.database.enablePostgres will make sure postgres service is configured and running. However password must be set manually with, ALTER USER [username] WITH PASSWORD '[password]';
      '' ++
      optional (cfg.database.password != "") ''
        config.services.goshrt.database.password will be stored as plaintext in the Nix store. Use database.passwordFile instead (when it's implemented).
      '';


  };
}
