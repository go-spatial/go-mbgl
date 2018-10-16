package mbgl

/*
#include "map.h"
*/
import "C"

type Map C.MbglMap

func NewMap(frontend RendererFrontend,
	size Size,
	pixelRatio float32,
	src FileSource,
	sched Scheduler) *Map {

	ptr := C.mbgl_map_new(
		frontend.rendererFrontend(),
		size.cSize(),
		C.float(pixelRatio),
		src.fileSource(),
		sched.scheduler())

	return (*Map)(ptr)
}

func (m *Map) Destruct() {
	C.mbgl_map_destruct((*C.MbglMap)(m))
}

func (m *Map) jumpTo(opts *CameraOptions) {
	C.mbgl_map_jump_to(
		(*C.MbglMap)(m),
		opts.cPtr())
}

func (m *Map) setStyleUrl(str string) {
	C.mbgl_map_set_style_url(
		(*C.MbglMap)(m),
		C.CString(str))
}
