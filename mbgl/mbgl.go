package mbgl

/*
#cgo CFLAGS: -fPIC
#cgo CFLAGS: -D_GLIBCXX_USE_CXX11_ABI=1
#cgo CXXFLAGS: -std=c++14 -std=gnu++14
#cgo CXXFLAGS: -g
#cgo CXXFLAGS: -I${SRCDIR}/../mason_packages/.link/include
#cgo CXXFLAGS: -I${SRCDIR}/../mapbox-gl-native/include
#cgo CXXFLAGS: -I${SRCDIR}/../mapbox-gl-native/platform/default
#cgo LDFLAGS: -L${SRCDIR}/mason_packages/.link/lib
#cgo LDFLAGS: -lsqlite3 -lz
#cgo LDFLAGS: -lmbgl-filesource -lmbgl-core

*/
import "C"
