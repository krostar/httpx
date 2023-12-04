{pkgs}:
pkgs.mkShellNoCC {
  nativeBuildInputs = with pkgs; [
    act
    deadnix
    alejandra
    gci
    git
    go_1_21
    gofumpt
    golangci-lint
    govulncheck
    shellcheck
    statix
    yamllint
  ];
}
