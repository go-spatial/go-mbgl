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
	"fmt"
)

type MapSnapshotter C.MbglMapSnapshotter


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

type PremultipliedImage C.MbglPremultipliedImage

func (im *PremultipliedImage) premultipliedImage() *C.MbglPremultipliedImage {
	return (*C.MbglPremultipliedImage)(im)
}

func (im *PremultipliedImage) Image() image.Image {
	raw := C.mbgl_premultiplied_image_raw(im.premultipliedImage())

	bytes := int(raw.width) * int(raw.height) * 4

	return &img{
		data: C.GoBytes(unsafe.Pointer(raw.data), C.int(bytes)),
		width: int(raw.width),
		height: int(raw.height),
	}

}

func (im *PremultipliedImage) Destruct() {
	C.mbgl_premultiplied_image_destruct(im.premultipliedImage())
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