{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
	nativeBuildInputs = [
		pkgs.go_1_21
		pkgs.zip
	];
}
