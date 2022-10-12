{ buildGoModule
, nix-gitignore
}:

let
  nixpkgs = import (fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/22.05.tar.gz";
    sha256 = "sha256:0d643wp3l77hv2pmg2fi7vyxn4rwy0iyr8djcw1h5x72315ck9ik";
  }) {
    config = { };
    overlays = [ ];
  };
in buildGoModule.override { go = nixpkgs.go_1_18; } rec {
  pname = "beagle";
  version = "1.1.3";

  src = nix-gitignore.gitignoreSource [ ] ./.;

  # The checksum of the Go module dependencies. `vendorSha256` will change if go.mod changes.
  # If you don't know the hash, the first time, set:
  # sha256 = "0000000000000000000000000000000000000000000000000000";
  # then nix will fail the build with such an error message:
  # hash mismatch in fixed-output derivation '/nix/store/m1ga09c0z1a6n7rj8ky3s31dpgalsn0n-source':
  # wanted: sha256:0000000000000000000000000000000000000000000000000000
  # got:    sha256:173gxk0ymiw94glyjzjizp8bv8g72gwkjhacigd1an09jshdrjb4
  vendorSha256 = "sha256-ikJtCG7DuN+dTYISbhGiY2JJFgDj1t8oWan8gXKa7NY=";
}
