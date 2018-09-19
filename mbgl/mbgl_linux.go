package mbgl

/*
#cgo LDFLAGS: -L${SRCDIR}/c/lib/linux
#cgo LDFLAGS: -lmbgl-loop-uv
#cgo LDFLAGS: -lsqlite -lnunicode -licu
#cgo LDFLAGS: -luv -lpthread -ldl -lcurl
#cgo LDFLAGS: -lpng16 -ljpeg
#cgo LDFLAGS: -lOSMesa
#cgo LDFLAGS: -lGL
*/
import "C"

/*
#cgo LDFLAGS: -lpng16 -ljpeg -lwebp
#cgo -licuuc -ldl
#cgo -lnu -lm
#cgo LDFLAGS: -lX11
#cgo LDFLAGS: -lOSMesa
#cgo LDFLAGS: -lGL
#cgo LDFLAGS: -lOpenGL -lGLX -lEGL
*/
