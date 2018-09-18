package mbgl

/*
#include "headless_frontend.h"
*/
import "C"

type HeadlessFrontend C.MbglHeadlessFrontend

func NewHeadlessFrontend(size *Size,
	pixelRatio float32,
	src FileSource,
	sched Scheduler,
	cacheDir *string,
	fontFamily *string) *HeadlessFrontend {

	var _cacheDir *C.char
	if cacheDir != nil {
		_cacheDir = C.CString(*cacheDir)
	}

	var _fontFamily *C.char
	if fontFamily != nil {
		_fontFamily = C.CString(*fontFamily)
	}

	ptr := C.mbgl_headless_frontend_new(
		size.cSize(),
		C.float(pixelRatio),
		src.fileSource(),
		sched.scheduler(),
		_cacheDir,
		_fontFamily)

	return (*HeadlessFrontend)(ptr)
}

func (hfe *HeadlessFrontend) rendererFrontend() *C.MbglRendererFrontend {
	return (*C.MbglRendererFrontend)(hfe)
}

func (hfe *HeadlessFrontend) Render(m *Map) *PremultipliedImage {
	ptr := C.mbgl_headless_frontend_render(
		(*C.MbglHeadlessFrontend)(hfe),
		(*C.MbglMap)(m))

	return (*PremultipliedImage)(ptr)
}

func (hfe *HeadlessFrontend) Destruct() {
	C.mbgl_headless_frontend_destruct((*C.MbglHeadlessFrontend)(hfe))
}
