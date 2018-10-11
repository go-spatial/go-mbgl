package simplified

/*
#cgo CFLAGS: -fPIC
#cgo CFLAGS: -D_GLIBCXX_USE_CXX11_ABI=1
#cgo CXXFLAGS: -std=c++14 -std=gnu++14
#cgo CXXFLAGS: -g
*/
import "C"
import "github.com/go-spatial/go-mbgl/mbgl"

var RunLoop *mbgl.RunLoop

func NewRunLoop() {
	if RunLoop == nil {
		RunLoop = mbgl.NewRunLoop()
	}
}

func DestroyRunLoop() {
	if RunLoop != nil {
		RunLoop.Destruct()
		RunLoop = nil
	}
}
