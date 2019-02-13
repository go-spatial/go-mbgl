# go-mbgl

go-mbgl aims to provide Go bindings for [Mapbox GL native's C++ API.](https://github.com/mapbox/mapbox-gl-native)

**WARNING**: This project is under heavly development, and is not production ready. 

**WARNING**: This will only work on Linux

This repository depends on [*git-lfs*](https://git-lfs.github.com/)

## Repoistory Layout

* [cmd/snap](cmd/snap) -- Holds the primary CLI (command line interface) tool. This tool servers both as a raster tile server, and static map generator.
* [mbgl/simplifed](mbgl/simplifed) -- Go bindings for snapshot API.
* [mbgl/c](mbgl/c) -- the C bridge for Go. (_Not working?_)
  * [mbgl/c/linux/lib](mbgl/c/linux/lib) -- location of the precompiled Linux libraries. 


## Development Env

Running the following commands should set up a working development environment

```console

docker build -t mbgl .

docker run -it -v it"$(pwd)":/go/src/github.com/go-spatial/go-mbgl

```

## Build snap commandline utility.

```console

cd cmd/snap

go build

```




