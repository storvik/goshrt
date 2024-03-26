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
  forEachHost = genAttrs cfg.nginx.hostnames;
  virtualHostsConfig = forEachHost
    (name: {
      forceSSL = true;
      enableACME = true;
      locations."/" = {
        extraConfig = ''
          proxy_max_temp_file_size 0;
          proxy_set_header Host               $host;
          proxy_set_header X-Real-IP          $remote_addr;
          proxy_set_header X-Forwarded-For    $proxy_add_x_forwarded_for;
          proxy_set_header X-Forwarded-Proto  https;
        '';
        proxyPass = "http://goshrt";
      };
    } // cfg.nginx.extraConfig);
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
    nginx = {
      enable = mkOption {
        default = false;
        type = types.bool;
        description = lib.mdDoc ''
          Enable nginx proxy. If not enabled proxy must be set up manually.
          Use extraArgs to set up reverse proxy / pass other options to
          services.nginx.
        '';
      };
      hostnames = mkOption {
        default = [ ];
        type = types.listOf types.str;
        description = lib.mdDoc "List of hostnames to work with goshrt.";
      };
      extraConfig = mkOption {
        default = { };
        type = types.anything;
        description = lib.mdDoc ''
          Attribute set with extra nginx config. For example to enable SSL
          and ACME add `{ forceSSL = true; enableACME = true; }`.
        '';
      };
    };
    database = {
      enable = mkOption {
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

    services.postgresql = mkIf cfg.database.enable {
      enable = true;
      ensureDatabases = [ cfg.database.name ];
      ensureUsers = [
        {
          name = cfg.database.user;
          ensureDBOwnership = true;
        }
      ];
    };

    services.nginx = mkIf cfg.nginx.enable {
      enable = true;
      upstreams.goshrt = {
        servers = { "localhost:${(toString cfg.httpPort)}" = { }; };
      };
      virtualHosts = virtualHostsConfig;
    };

    warnings =
      optional cfg.database.enable ''
        config.services.goshrt.database.enable will make sure postgres service is configured and running. However password must be set manually with, ALTER USER [username] WITH PASSWORD '[password]';
      '' ++
      optional (cfg.database.password != "") ''
        config.services.goshrt.database.password will be stored as plaintext in the Nix store. Use database.passwordFile instead (when it's implemented).
      '';


  };
}
