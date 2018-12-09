# MiSTer-bootstrap

MiSTer bootstrap is a tool that updates all cores for the MiSTer platform. Can be compiled to a binary for multiple platforms, and also a shared library to be directly invoked from C-code.

# Requirements

* Go
* Godep

# Usage

`./bootstrap`
* `-r <Repo URL>`
	* Optional Repo URL, defaults to OpenVGS/MiSTer-repository which is updated hourly.
* `-o <Output Directory>`
	* Optional Output Directory, defaults to current directory.
