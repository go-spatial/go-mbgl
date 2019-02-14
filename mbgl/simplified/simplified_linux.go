package simplified

/*
#cgo CXXFLAGS: -I${SRCDIR}/../c/include
#cgo CXXFLAGS: -I${SRCDIR}/../c/include/include
#cgo LDFLAGS: -L${SRCDIR}/../c/lib/linux

#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libmbgl-core.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libmbgl-filesource.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libmbgl-loop-uv.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libsqlite.a

#cgo LDFLAGS: -Wl,--no-as-needed -lcurl
#cgo LDFLAGS: -Wl,--as-needed ${SRCDIR}/../c/lib/linux/libmbgl-core.a
#cgo LDFLAGS: -lOSMesa
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libpng.a
#cgo LDFLAGS: -lz -lm
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libjpeg.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libnunicode.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libicu.a
#cgo LDFLAGS: ${SRCDIR}/../c/lib/linux/libuv.a
#cgo LDFLAGS: -lrt -lpthread
#cgo LDFLAGS: -lnsl -ldl
#cgo LDFLAGS: -static-libstdc++

*/
import "C"

// @cgo CXXFLAGS: -I${SRCDIR}/../mason_packages/.link/include
// @cgo CXXFLAGS: -I${SRCDIR}/../mapbox-gl-native/vendor/expected/include
// @cgo CXXFLAGS: -I${SRCDIR}/../mapbox-gl-native/include
// @cgo CXXFLAGS: -I${SRCDIR}/../mapbox-gl-native/platform/default
// @cgo LDFLAGS: -L${SRCDIR}/../mason_packages/.link/lib
// @cgo LDFLAGS: -lsqlite3 -lz
// @cgo LDFLAGS: -lmbgl-filesource -lmbgl-core
// #cgo CXXFLAGS: -I${SRCDIR}/c/mapbox-gl-native/include
// #cgo LDFLAGS: -Wl,-Bsymbolic-functions ${SRCDIR}/../c/lib/linux/libuv.a
