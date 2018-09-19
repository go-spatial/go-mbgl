# PLACE YOUR MAPBOXGL INCLUDES HERE
INCLUDES=\
          	-I/Users/julio/lib/c/mapbox-gl-native/include \
          	-I/Users/julio/lib/c/mapbox-gl-native/platform/default

LINKS_DW=\
	-lmbgl-loop-darwin\
	-framework Mapbox\
	-framework CoreFoundation\
	-framework CoreGraphics\
	-framework ImageIO\
	-framework OpenGL\
	-framework CoreText\
	-framework Foundation

# reflecetd in mblg.go file
BUILD_FLAGS=-fPIC\
                -D_GLIBCXX_USE_CXX11_ABI=1\
                -std=c++14 -std=gnu++14\
                -g\
                -I./mason_packages/\.link/include\
                -L./mason_packages/.link/lib\
                -lsqlite3 -lz\
                -lmbgl-filesource -lmbgl-core

CGO_CXXFLAGS=$(INCLUDES)

GO_FLAGS=CGO_CXXFLAGS="$(CGO_CXXFLAGS)"

run=
dir=./...

.PHONY: test
test:
	$(GO_FLAGS) go test -v $(dir) --run=$(run)

f=
.PHONY: dwfile
dwfile:
	g++ $(BUILD_FLAGS) $(CGO_CXXFLAGS) $(LINKS_DW) -c $(f)
