{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    buf
    protoc-gen-go
    protoc-gen-connect-go
    protoc-gen-es
    go
    grpcui
    nodejs_24
    pnpm
  ];
}
