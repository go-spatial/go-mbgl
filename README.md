# Mbgl
go bindings for mapbox gl native's c++ api


1. Download and build [Mapbox-GL-Native](https://github.com/mapbox/mapbox-gl-native) for your platform.

2. Add the absolute or relative path to the `include` and `platform/default` to the CGO_CXXFLAGS variable in the makefile as `-I` flags
  
3. Make library binaries accessible to linker
    * For `darwin` this requires copying a `.framework` file into `/Libray/Frameworks` and `.a` files into `./mason_packages/.link/lib`
    * For `linux` this can be done by copying `.a` files into `./mason_packages/.link/lib`

4. Install [mason-js](https://github.com/mapbox/mason-js) and run `mason-js install` and then `mason-js link` in the atto root directory to download the required dependencies 

5. `go build` to build...

There is a Dockerfile in /docker which details the steps required to set up a build and runtime environment for go-mbgl.