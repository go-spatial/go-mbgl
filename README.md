# go-mbgl

go-mbgl aims to provide Go bindings for [Mapbox GL native's C++ API.](https://github.com/mapbox/mapbox-gl-native)

**WARNING**: This project is under heavly development, and is not production ready. 

**WARNING**: This will only work on Linux

This repository depends on [*git-lfs*](https://git-lfs.github.com/)

## Repoistory Layout

* [cmd/snap](cmd/snap) -- Holds the primary CLI (command line interface) tool. This tool servers both as a raster tile server, and static map generator.
* [mbgl/simplifed](mbgl/simplifed) -- Go bindings for snapshot API.
* [mbgl/c](mbgl/c) -- the C bridge for Go.
  * [mbgl/c/linux/lib](mbgl/c/linux/lib) -- location of the precompiled Linux libraries. 

