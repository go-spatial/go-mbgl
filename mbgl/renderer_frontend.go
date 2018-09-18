package mbgl

/*
#include "renderer_frontend.h"
*/
import "C"

type RendererFrontend interface {
	rendererFrontend() *C.MbglRendererFrontend
	Destruct()
}
