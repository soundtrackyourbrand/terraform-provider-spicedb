with (import <nixpkgs> {config.allowUnfree = true;});

{
  shell = mkShell {
    buildInputs = [ terraform go gotools gopls ];

    shellHook = ''
      export GOPATH=$(go env GOPATH)
    '';
  };
}

