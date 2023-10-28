with import <nixpkgs> { };

{
  shell = mkShell {
    buildInputs = [ terraform go gotools gopls ];

    shellHook = ''
      export GOPATH=$(go env GOPATH)
    '';
  };
}

