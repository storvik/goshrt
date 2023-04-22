{ pkgs, ... }:

{
  # Only allow this to boot as a container
  boot.isContainer = true;

  # Allow nginx through the firewall
  networking.firewall.allowedTCPPorts = [ 80 ];

  # Needed because of `nginx.virtualHosts..enableACME`
  security.acme.acceptTerms = true;
  security.acme.defaults.email = "goshrt@example.com";

  services.goshrt = {
    enable = true;
    httpPort = 8080;
    key = "qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$";
    database = {
      enable = true;
      host = "localhost";
      port = 5432;
      user = "goshrt";
      password = "trhsog";
    };
    nginx = {
      enable = true;
      hostnames = [ "examplename1.com" "examplename2.com" ];
      extraConfig = { forceSSL = true; enableACME = true; };
    };
  };

}
