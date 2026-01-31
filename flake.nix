{
  description = "Flake for datastarui";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "datastarui";
          version = "0.1.0";
          src = ./.;
          vendorHash = null;

          nativeBuildInputs = with pkgs; [
            templ
            tailwindcss
          ];
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            gotools
            go-tools
            templ
            air
            pnpm
            tailwindcss
            just
          ];

          shellHook = ''
            echo "DatastarUI Development Environment"
            echo ""
            echo "Commands:"
            echo "  just tailwind    - start the Tailwind CSS watcher"
            echo "  just watch       - start the Go server with live reload"
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/compline";
        };
      }
    );
}
