{
  inputs = {
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {nixpkgs-unstable, ...}: let
    supportedSystems = ["aarch64-linux" "aarch64-darwin" "x86_64-linux" "x86_64-darwin"];
    forEachSupportedSystems = f: builtins.listToAttrs (builtins.map (system: (nixpkgs-unstable.lib.nameValuePair system (f system))) supportedSystems);
    pkgsForSystem = system: {
      nixpkgs ? nixpkgs-unstable,
      overlays ? [],
    }: (import nixpkgs {inherit system overlays;});
  in {
    devShells = forEachSupportedSystems (system: {default = import ./shell.nix {pkgs = pkgsForSystem system {};};});
    formatter = forEachSupportedSystems (system: let pkgs = pkgsForSystem system {}; in pkgs.alejandra);
  };
}
