package mbgl

/*
#include "map_snapshotter.h"
*/
import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

type PremultipliedImage C.MbglPremultipliedImage

func (im *PremultipliedImage) premultipliedImage() *C.MbglPremultipliedImage {
	return (*C.MbglPremultipliedImage)(im)
}

func (im *PremultipliedImage) Image() image.Image {
	raw := C.mbgl_premultiplied_image_raw(im.premultipliedImage())

	bytes := int(raw.width) * int(raw.height) * 4

	return &img{
		data:   C.GoBytes(unsafe.Pointer(raw.data), C.int(bytes)),
		width:  int(raw.width),
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
