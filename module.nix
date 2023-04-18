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
    database = {
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
      # TODO: prestart should possibly run database migrations, but must add migrate command first

      serviceConfig = {
        Type = "simple";
        ExecStart = "${goshrt}";
        Restart = "always";
      };
    };

    warnings =
      optional (cfg.database.password != "") ''
        config.services.goshrt.database.password will be stored as plaintext in the Nix store. Use database.passwordFile instead (when it's implemented).
      '';


  };
}
