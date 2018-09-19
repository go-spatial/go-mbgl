// +build !linux

package mbgl

/*
#include "map_snapshotter.h"
*/
import "C"
import (
	"fmt"
	"strings"
)

type MapSnapshotter C.MbglMapSnapshotter

func NewMapSnapshotter(src FileSource,
	sched Scheduler,
	style string,
	size Size,
	pixelRatio float32,
	camOpts *CameraOptions,
	region *LatLngBounds,
	cacheDir *string) *MapSnapshotter {

	var isFile = 0
	if !strings.Contains(style, "http") {
		isFile = 1
	}

	var _camOpts *C.MbglCameraOptions
	if camOpts != nil {
		_camOpts = camOpts.cPtr()
	}

	var _region *C.MbglLatLngBounds
	if region != nil {
		_region = region.latLngBounds()
	}

	var _cacheDir *C.char
	if cacheDir != nil {
		_cacheDir = C.CString(*cacheDir)
	}

	ptr := C.mbgl_map_snapshotter_new(
		src.fileSource(),
		sched.scheduler(),
		C.int(isFile), C.CString(style),
		size.cSize(),
		C.float(pixelRatio),
		_camOpts,
		_region,
		_cacheDir)

	return (*MapSnapshotter)(ptr)
}

func (ms *MapSnapshotter) mapSnapshotter() *C.MbglMapSnapshotter {
	return (*C.MbglMapSnapshotter)(ms)
}

func (ms *MapSnapshotter) Snapshot() *PremultipliedImage {
	ptr := C.mbgl_map_snapshotter_snapshot(ms.mapSnapshotter())
	return (*PremultipliedImage)(ptr)
}

func (ms *MapSnapshotter) SetCameraOptions(camOpts CameraOptions) {
	C.mbgl_map_snapshotter_set_camera_options(ms.mapSnapshotter(), camOpts.cPtr())
}

func (ms *MapSnapshotter) SetRegion(region *LatLngBounds) {
	C.mbgl_map_snapshotter_set_region(ms.mapSnapshotter(), region.latLngBounds())
}

func (ms *MapSnapshotter) SetStyleURL(style string) {
	fmt.Println("setting style, ", style)
	C.mbgl_map_snapshotter_set_style_url(ms.mapSnapshotter(), C.CString(style))
	fmt.Println("style set")
}

func (ms *MapSnapshotter) SetSize(size Size) {
	C.mbgl_map_snapshotter_set_size(ms.mapSnapshotter(), size.cSize())
}

func (ms *MapSnapshotter) Destruct() {
	C.mbgl_map_snapshotter_destruct(ms.mapSnapshotter())
}
