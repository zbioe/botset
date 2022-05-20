{ pkgs, ... }:
{
  config.services = {

#    webapp = {
#      service = {
#        image = "stevepryde/thirtyfour_testapp:0.1.0";
#        ports = [
#          "8000:80"
#        ];
#      };
#    };
    selenium = {
      composition = {
        shm_size = "2gb";
      };
      service = {
        image = "selenium/standalone-chrome:4.1.0-20211123";
        ports = [
          "4444:4444"
          "7900:7900"
        ];
        environment = {
          NODE_MAX_INSTANCES = "4";
          NODE_MAX_SESSION = "4";
          SCREEN_WIDTH = "1360";
          SCREEN_HEIGHT = "1020";
          SCREEN_DEPTH = "24";
          SCREEN_DPI = "96";
          SE_NODE_MAX_SESSIONS = "4";
        };
      };
    };
  };
}
