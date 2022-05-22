{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # pkgs.systemfd
    vgo2nix
    go
    sqlite
  ];

  shellHook = ''
    unset GOPATH
  '';
}
