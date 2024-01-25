{
  description = "hosting.de terraform provider";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  inputs.nixpkgs-2305.url = "github:nixos/nixpkgs/nixos-23.05";
  inputs.devshell.url = "github:numtide/devshell";
  inputs.flake-parts.url = "github:hercules-ci/flake-parts";

  outputs = inputs@{ self, flake-parts, devshell, nixpkgs, nixpkgs-2305 }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        devshell.flakeModule
      ];

      systems = [
        "x86_64-linux"
      ];

      perSystem = { pkgs, ... }: {
        devshells.default = {
          # Add additional packages you'd like to be available in your devshell
          # PATH here
          packages = with pkgs; [
            go
            errcheck
            go-tools
            gnumake
            golangci-lint
            inputs.nixpkgs-2305.legacyPackages.${pkgs.system}.terraform
          ];
          bash.extra = ''
            export GOPATH=~/.local/share/go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      };
    };
}
