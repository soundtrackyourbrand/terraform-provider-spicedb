with import <nixpkgs> { };

{
  shell = mkShell {
    buildInputs = [ terraform go gotools pcre pcre2];

    shellHook = ''
      export GOPATH=$(go env GOPATH)
    '';
  };
}

