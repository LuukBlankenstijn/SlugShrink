{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/25.11";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [ buf go grpcui nodejs_24 pnpm ];
        shellHook = ''
          export COREPACK_ENABLE_DOWNLOADS=1
          export COREPACK_HOME="$PWD/.corepack"
          export PATH="$COREPACK_HOME/shims:$PATH"
          corepack prepare pnpm@latest --activate >/dev/null 2>&1 || true
        '';
      };
    };
}
