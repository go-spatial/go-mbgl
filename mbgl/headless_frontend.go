package mbgl

/*
#include "headless_frontend.h"
*/
import "C"

import (
	"fmt"
	"os"
)

type HeadlessFrontend C.MbglHeadlessFrontend

func NewHeadlessFrontend(size *Size, pixelRatio float32, src FileSource, sched Scheduler, cacheDir *string, fontFamily *string) *HeadlessFrontend {

	var _cacheDir *C.char
	if cacheDir != nil {
		_cacheDir = C.CString(*cacheDir)
	}

	var _fontFamily *C.char
	if fontFamily != nil {
		_fontFamily = C.CString(*fontFamily)
	}
	fmt.Fprintf(os.Stderr, "Starting to startup mbgl_headless_frontend_new\n")

	ptr := C.mbgl_headless_frontend_new(
		size.cSize(),
		C.float(pixelRatio),
		src.fileSource(),
		sched.scheduler(),
		_cacheDir,
		_fontFamily)
	fmt.Fprintf(os.Stderr, "After call to mbgl_headless_frontend_new")

	return (*HeadlessFrontend)(ptr)
}

func (hfe *HeadlessFrontend) rendererFrontend() *C.MbglRendererFrontend {
	return (*C.MbglRendererFrontend)(hfe)
}

func (hfe *HeadlessFrontend) Render(m *Map) *PremultipliedImage {
	fmt.Fprintf(os.Stderr, "Before the Render")
	ptr := C.mbgl_headless_frontend_render(
		(*C.MbglHeadlessFrontend)(hfe),
		(*C.MbglMap)(m))
	fmt.Fprintf(os.Stderr, "After the Render")

	return (*PremultipliedImage)(ptr)
}

func (hfe *HeadlessFrontend) Destruct() {
	C.mbgl_headless_frontend_destruct((*C.MbglHeadlessFrontend)(hfe))
}
