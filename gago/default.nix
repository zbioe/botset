{ buildGoPackage, lib }:

buildGoPackage {
  pname = "joca";
  version = "0.1";

  goPackagePath = "github.com/zbioe/joca";

  src = lib.cleanSource ./.;

  goDeps = ./deps.nix;

  meta = {
    description = "Bot twitter";
    homepage = https://github.com/zbioe/joca;
    license = lib.licenses.wtfpl;
    maintainers =  [ lib.maintainers.zbioe ];
  };
}
