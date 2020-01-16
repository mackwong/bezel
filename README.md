# Bezel

This project is the prototype of diamond-on-edge cli tool to configure cluster shaping and distrbute to devices.

[![build status](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/badges/master/build.svg)](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/commits/master)
[![coverage report](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/badges/master/coverage.svg)](http://gitlab.bj.sensetime.com/diamond/service-providers/bezel/commits/master)

## Build

`make arm64/amd64/windows`

Support building executable binary on different architectures, include arm64(arch64), amd64 and windows.

## Usage

`bezel-arm64 create`

To generate global config file, `edge-config.yaml`, and sub config files in `${pwd}/sub`.

After binary executing, you should give inputs with interactive command line.

Follow the hints on screen, which would print the fields you need configure.

## Binary 

Coming soon