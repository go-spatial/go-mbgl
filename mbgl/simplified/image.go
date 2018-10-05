package simplified

import (
	"image"
	"image/color"
)

type Image struct {
	Data          []byte
	Width, Height int
}

func (im Image) ColorModel() color.Model { return color.RGBAModel }
func (im Image) Bounds() image.Rectangle { return image.Rect(0, 0, im.Width, im.Height) }
func (im Image) At(x, y int) color.Color {
	i := im.Width * 4 * y
	i += 4 * x

	return color.RGBA{
		im.Data[i],
		im.Data[i+1],
		im.Data[i+2],
		im.Data[i+3],
	}
}
