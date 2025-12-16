{
  description = "HTTP web server for SPA";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-utils.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachSystem
      [ "x86_64-linux" ]
      (
        system:
        let
          nixpkgs' = nixpkgs.legacyPackages.${system};
          pkgs = import self {
            inherit system;
            inputs = null;
            nixpkgs = nixpkgs';
          };
        in
        {
          defaultPackage = pkgs.default;
          packages = pkgs;
          devShells.default = pkgs.devShell;
          checks = {
            fmt =
              with nixpkgs';
              runCommandLocal "fmt" { } ''
                export HOME=$(mktemp -d)
                cp -r ${./.} src
                chmod -R u+w src
                cd src
                export PATH=${
                  lib.makeBinPath [
                    nixfmt-rfc-style
                    gofumpt
                  ]
                }:$PATH
                ${treefmt}/bin/treefmt --fail-on-change > $out
              '';
          };
        }
      );
}
