{
  description = "hosting.de terraform provider";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    devshell.url = "github:numtide/devshell";
    devshell.inputs.flake-utils.follows = "flake-utils";
    devshell.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, flake-utils, devshell, nixpkgs }:
    flake-utils.lib.simpleFlake {
      inherit self nixpkgs;
      name = "hosting.de terraform provider";
      overlay = devshell.overlays.default;
      shell = { pkgs }:
        pkgs.devshell.mkShell {
          # Add additional packages you'd like to be available in your devshell
          # PATH here
          devshell.packages = with pkgs; [
            go
            errcheck
            go-tools
            gnumake
            golangci-lint
          ];
          bash.extra = ''
            export GOPATH=~/.local/share/go
            export PATH=$GOPATH:$PATH
          '';
        };
    };
}
