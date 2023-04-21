{ pkgs, ... }:

{
  # Only allow this to boot as a container
  boot.isContainer = true;

  # Allow nginx through the firewall
  networking.firewall.allowedTCPPorts = [ 80 ];

  services.goshrt = {
    enable = true;
    httpPort = 8080;
    key = "qTGVn$a&hRJ9385C^z7L!MW5CnwZq3&$";
    database = {
      enablePostgres = true;
      host = "localhost";
      port = 5432;
      user = "goshrt";
      password = "trhsog";
    };
  };

}
