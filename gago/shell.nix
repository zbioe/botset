{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    # pkgs.systemfd
    vgo2nix
    go_1_17
    sqlite
  ];

  shellHook = ''
    unset GOPATH
  '';
}
