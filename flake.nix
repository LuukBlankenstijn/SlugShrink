{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/25.11";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          buf
          protoc-gen-go
          protoc-gen-connect-go
          go
          grpcui
        ];
      };
    };
}
