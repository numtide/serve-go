{
  system ? builtins.currentSystem,
  inputs ? import ./flake.lock.nix {},
  nixpkgs ?
    import inputs.nixpkgs {
      inherit system;
      # Makes the config pure as well. See <nixpkgs>/top-level/impure.nix:
      config = {};
      overlays = [];
    },
  buildGoModule ? nixpkgs.buildGoModule,
}: let
  serve-go =
    buildGoModule
    {
      name = "serve-go";
      src = ./.;
      vendorSha256 = null;
      meta = with nixpkgs.lib; {
        description = "HTTP web server for SPA";
        homepage = "https://github.com/numtide/serve-go";
        license = licenses.mit;
        maintainers = with maintainers; [zimbatm jfroche];
        platforms = platforms.linux;
      };
    };
  devShell =
    nixpkgs.mkShellNoCC
    {
      buildInputs = with nixpkgs; [
        gofumpt
        golangci-lint
        alejandra
        go
        golint
        treefmt
        just
        gcc
      ];
    };
in {
  inherit serve-go devShell;
  default = serve-go;
}
