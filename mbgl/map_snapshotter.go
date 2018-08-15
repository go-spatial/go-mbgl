package mbgl

/*
#include "map_snapshotter.h"
 */
import "C"
import (
	"strings"
	"image"
	"image/color"
	"unsafe"
)

type MapSnapshotter struct {
	ptr *C.MbglMapSnapshotter
}


func NewMapSnapshotter(src FileSource,
	sched Scheduler,
	style string,
	size Size,
	pixelRatio float32,
	camOpts *CameraOptions,
	region *LatLngBounds,
	cacheDir *string)*MapSnapshotter {

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
			_region = region.ptr
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

		return &MapSnapshotter{ptr:ptr}
}

func (ms MapSnapshotter) Snapshot() *PremultipliedImage {
	ptr := C.mbgl_map_snapshotter_snapshot(ms.ptr)
	return &PremultipliedImage{ptr:ptr}
}

func (ms MapSnapshotter) SetCameraOptions(camOpts CameraOptions) {
	C.mbgl_map_snapshotter_set_camera_options(ms.ptr, camOpts.cPtr())
}

func (ms MapSnapshotter) SetRegion(region *LatLngBounds) {
	C.mbgl_map_snapshotter_set_region(ms.ptr, region.ptr)
}

func (ms *MapSnapshotter) Destruct() {
	C.mbgl_map_snapshotter_destruct(ms.ptr)
	ms.ptr = nil
}

type PremultipliedImage struct {
	ptr *C.MbglPremultipliedImage
}

func (im PremultipliedImage) Image() image.Image {
	raw := C.mbgl_premultiplied_image_raw(im.ptr)

	bytes := int(raw.width) * int(raw.height) * 4

	return &img{
		data: C.GoBytes(unsafe.Pointer(raw.data), C.int(bytes)),
		width: int(raw.width),
		height: int(raw.height),
	}

}

type img struct {
	data          []byte
	width, height int
}

func (im img) ColorModel() color.Model {
	return color.RGBAModel
}

func (im img) Bounds() image.Rectangle {
	return image.Rect(0, 0, im.width, im.height)
}

func (im img) At(x, y int) color.Color {
	i := im.width * 4 * y
	i += 4 * x

	return color.RGBA{
		im.data[i],
		im.data[i+1],
		im.data[i+2],
		im.data[i+3],
	}
}