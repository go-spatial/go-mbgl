# Mbgl
go bindings for mapbox gl native's c++ api

## Development

### Dev Deps
**Darwin**:
see inssructions [on the mblg page](https://github.com/mapbox/mapbox-gl-native/blob/master/platform/macos/INSTALL.md), but do not build the library.
In list form:
* xcode
* node
* cmake
* ccache
* xcpretty
* jazzy

**Linux**:
todo

### install
run the install script provided in the repo. Note, it'll install a node package, [`mason-js`](https://github.com/mapbox/mason-js) globally with `npm i -g`.
```bash
./install.sh
```

### build/tests
regular go tool commands (`go build`, `go test`) can be used for building the library.

### Old install instructions (deprecated)

1. Download and build [Mapbox-GL-Native](https://github.com/mapbox/mapbox-gl-native) for your platform.

2. Add the absolute or relative path to the `include` and `platform/default` to the CGO_CXXFLAGS variable in the makefile as `-I` flags
  
3. Make library binaries accessible to linker
    * For `darwin` this requires copying a `.framework` file into `/Libray/Frameworks` and `.a` files into `./mason_packages/.link/lib`
    * For `linux` this can be done by copying `.a` files into `./mason_packages/.link/lib`

4. Install [mason-js](https://github.com/mapbox/mason-js) and run `mason-js install` and then `mason-js link` in the atto root directory to download the required dependencies 

5. `go build` to build...

There is a Dockerfile in /docker which details the steps required to set up a build and runtime environment for go-mbgl.