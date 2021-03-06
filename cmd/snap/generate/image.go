package generate

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/go-mbgl/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
)

type ReadAtCloser interface {
	io.ReaderAt
	io.Closer
}

type CenterRect struct {
	Lat  float64
	Lng  float64
	Rect image.Rectangle
	// for backing store
	offset   int64
	length   int
	imgWidth int
}

type Image struct {
	// Width of the desired image, it will be multipiled by the PPIRatio to get the final width
	width int
	// Height of the desired image, it will be multipiled by the PPIRatio to get the final height.
	height int

	// These are the centers and the rectangles where the image will be
	// placed
	centers []CenterRect

	// the offset from the top, this is for clip
	offsetHeight int
	offsetWidth  int

	// PPIRatio
	ppiratio float64
	// Style to use to generate the tile
	style string
	// The zoom level
	zoom float64

	// Projection
	prj bounds.AProjection

	// We will write the data to this and then use this for the
	// At function.
	backingStore *os.File
	initilized   bool

	// this is for debugging.
	drawBounds bool
	// bounds need to be in the coordinate system of the image.
	// the color will be black
	bounds     [4]float64
	fullBounds image.Rectangle
}

func (im *Image) SetDebugBounds(extent *geom.Extent, zoom float64) {

	const tilesize = 4096 / 2
	const scale = 4

	// for lat lng geom.Extent should be laid out as follows:
	// {west, south, east, north}
	ne := [2]float64{extent[3], extent[2]}
	sw := [2]float64{extent[1], extent[0]}

	swPt := bounds.LatLngToPoint(im.prj, sw[0], sw[1], zoom, tilesize)
	nePt := bounds.LatLngToPoint(im.prj, ne[0], ne[1], zoom, tilesize)
	im.drawBounds = true
	im.bounds = [4]float64{
		float64(int((nePt[0] - float64(im.fullBounds.Min.X)) / scale)),
		float64(int((swPt[1] - float64(im.fullBounds.Min.Y)) / scale)),
		float64(int((swPt[0] - float64(im.fullBounds.Min.X)) / scale)),
		float64(int((nePt[1] - float64(im.fullBounds.Min.Y)) / scale)),
	}
}

func (_ Image) ColorModel() color.Model { return color.RGBAModel }
func (im Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(float64(im.width)*im.ppiratio), int(float64(im.height)*im.ppiratio))
}

func (im Image) Close() {
	if im.backingStore == nil {
		return
	}
	log.Println("Closing backing store file:", im.backingStore.Name())
	im.backingStore.Close()
	// ignore any errors.
	_ = os.Remove(im.backingStore.Name())
}

func (im Image) At(x, y int) color.Color {
	rx, ry := x+im.offsetWidth, y+im.offsetHeight
	// rx, ry := x, y
	var data [4]byte

	if im.drawBounds {
		if int(im.bounds[0]) == rx || int(im.bounds[2]) == rx || int(im.bounds[1]) == ry || int(im.bounds[3]) == ry {
			return color.RGBA{0, 0, 0, 255}
		}
	}

	// We need to look through the centers to find the first rect that containts this x,y
	for i := range im.centers {
		rect := im.centers[i].Rect
		if rect.Min.X <= rx && rx <= rect.Max.X && rect.Min.Y <= ry && ry <= rect.Max.Y {
			dx, dy := rx-rect.Min.X, ry-rect.Min.Y
			idx := int64(im.centers[i].imgWidth*4*dy+(4*dx)) + (im.centers[i].offset)
			_, err := im.backingStore.ReadAt(data[:], idx)
			if err != nil {
				panic(fmt.Sprintf("(%v,%v) -> Centers[%v]{ %v }: %v Got an error reading backing store: %v", x, y, i, im.centers[i], idx, err))
			}
			return color.RGBA{data[0], data[1], data[2], data[3]}
		}
	}
	panic(fmt.Sprintf("Did not find expected offset %v,%v -- %v,%v", x, y, rx, ry))
	return color.RGBA{}
}

func NewImage(
	prj bounds.AProjection,
	desiredWidth, desiredHeight int,
	centerXY [2]float64,
	zoom float64,
	ppi, pitch, bearing float64,
	style string,
	tempDir, tempFilename string,
) (*Image, error) {

	const tilesize = 4096 / 2
	const scale = 4

	numTilesNeeded := int(
		math.Ceil((math.Max(
			float64(desiredWidth),
			float64(desiredHeight),
		)/tilesize + 1) / 2,
		),
	)
	centerTileLength := int(math.Ceil((tilesize - 1) * ppi))

	tmpDir := "."
	if tempDir == "" {
		tmpDir = tempDir
	}
	tmpFilename := "image_backingstore.bin."
	if tempFilename == "" {
		tmpFilename = tempFilename
	}

	tmpfile, err := ioutil.TempFile(tmpDir, tmpFilename)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup backing store: %v", err)
	}

	img := Image{
		prj:          prj,
		style:        style,
		zoom:         zoom,
		width:        desiredWidth,
		height:       desiredHeight,
		ppiratio:     ppi,
		centers:      make([]CenterRect, 0, numTilesNeeded*numTilesNeeded),
		backingStore: tmpfile,
	}

	ry := 0
	rx := 0
	bsOffset := int64(0)
	var min image.Point
	for y := -numTilesNeeded; y <= numTilesNeeded; y++ {
		rx = 0
		for x := -numTilesNeeded; x <= numTilesNeeded; x++ {

			var crect CenterRect
			center := [2]float64{centerXY[0] + (float64(x*tilesize) * scale), centerXY[1] + (float64(y*tilesize) * scale)}
			crect.Lat, crect.Lng = bounds.PointToLatLng(prj, center, zoom, tilesize)
			crect.Rect = image.Rect(rx, ry, rx+centerTileLength, ry+centerTileLength)
			// This is the top-left most coner
			if rx == 0 && ry == 0 {
				min.X, min.Y = int(center[0]-(tilesize/2*(scale))), int(center[1]-(tilesize/2*(scale)))
			}
			snpsht := mbgl.Snapshotter{
				Style:    style,
				Width:    uint32(tilesize),
				Height:   uint32(tilesize),
				PPIRatio: ppi,
				Lat:      crect.Lat,
				Lng:      crect.Lng,
				Zoom:     zoom,
			}
			snpImage, err := mbgl.Snapshot1(snpsht)
			if err != nil {
				// Delete the tempfile
				img.Close()
				return nil, err
			}
			crect.length, err = img.backingStore.Write(snpImage.Data)
			if err != nil {
				// Delete the tempfile
				img.Close()
				return nil, err
			}
			crect.offset = bsOffset
			crect.imgWidth = snpImage.Width
			img.centers = append(img.centers, crect)

			bsOffset += int64(crect.length)
			rx += centerTileLength
		}
		ry += centerTileLength
	}
	img.fullBounds.Min = min
	img.offsetWidth = (rx / 2) - int(float64(desiredWidth/2)*ppi)
	img.offsetHeight = (ry / 2) - int(float64(desiredHeight/2)*ppi)

	log.Println("Done generating images")
	img.initilized = true
	err = img.backingStore.Sync()
	// Move to the top of the file.
	log.Printf("Backing store has been sync'd : %v -- %v", img.backingStore.Name(), err)
	_, _ = img.backingStore.Seek(0, 0)
	return &img, nil

}
