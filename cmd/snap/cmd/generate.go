package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateWidth int
var cmdGenerateHeight int
var cmdGeneratePPIRatio float64
var cmdGeneratePitch float64
var cmdGenerateBearing float64
var cmdGenerateFormat string
var cmdGenerateOutputName string

var cmdGenerate = &cobra.Command{
	Use:     "generate",
	Short:   "generate an image for given map coordinates.",
	Aliases: []string{"gen"},
}

func ValidateGenerateParams() error {
	if cmdGenerateWidth < 1 {
		return errors.New("--width must be greater then or equal to 1")
	}
	if cmdGenerateHeight < 1 {
		return errors.New("--height must be greater then or equal to 1")
	}
	if cmdGeneratePitch < 0 {
		return errors.New("--pitch must be greater then or equal to 0")
	}
	if cmdGenerateBearing < 0 {
		return errors.New("--bearing must be greater then or equal to 0")
	}
	cmdGenerateFormat = strings.ToLower(cmdGenerateFormat)
	if cmdGenerateFormat != "" && (cmdGenerateFormat != "jpg" || cmdGenerateFormat != "png") {
		return errors.New("--format must be jpg or png")
	}
	return nil
}

func IsValidLngString(lng string) (float64, bool) {

	f64, err := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if err != nil {
		return f64, false
	}
	return f64, -180.0 <= f64 && f64 <= 190.0
}

func IsValidLatString(lat string) (float64, bool) {

	f64, err := strconv.ParseFloat(strings.TrimSpace(lat), 64)
	if err != nil {
		return f64, false
	}
	return f64, -90.0 <= f64 && f64 <= 90.0
}

func getOutfile(output string, fn func(w io.Writer, ext string) error) error {
	ext, err := getFormatString(output)
	if err != nil {
		return err
	}
	var out io.Writer

	if output == "" {
		out = os.Stdout
	} else {
		file, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("error creating output file: %v -- %v\n", output, err)
		}
		defer file.Close()
		out = file
	}
	return fn(out, ext)
}

func genImage(center [2]float64, zoom float64, output string) error {
	return getOutfile(output, func(out io.Writer, ext string) error {

		_, err := generateZoomCenterImage(
			out,
			RootStyle,
			cmdGenerateWidth, cmdGenerateHeight,
			cmdGeneratePPIRatio,
			cmdGeneratePitch, cmdGenerateBearing,
			center, zoom, ext,
		)
		if err != nil {
			return fmt.Errorf("failed to generate image: %v", err)
		}

		return nil
	})
}

func genCenterZoomTilesImage(centerPtX, centerPtY, zoom float64, tilesize int) (*image.NRGBA, error) {

	const padding = 10
	// assume ESPG3857 for now.
	prj := bounds.ESPG3857

	minStartingPtX := centerPtX - float64(cmdGenerateWidth/2)

	startingPtX := centerPtX

	for startingPtX >= minStartingPtX {
		startingPtX -= float64(tilesize)
	}

	startingPtY := centerPtY
	minStartingPtY := centerPtY - float64(cmdGenerateHeight/2)
	for startingPtY >= minStartingPtY {
		startingPtY -= float64(tilesize)
	}
	var centers [][2]float64
	var rects []image.Rectangle

	zeroPtX := int(startingPtX) - (tilesize+padding)/2
	zeroPtY := int(startingPtY) - (tilesize+padding)/2
	var maxWidth, maxHeight int

	{
		y := startingPtY
		for ; y <= centerPtY+float64((cmdGenerateHeight+padding)/2); y += float64(tilesize) {
			x := startingPtX

			for ; x <= centerPtX+float64((cmdGenerateWidth+padding)/2); x += float64(tilesize) {
				lat, lng := bounds.PointToLatLng(prj, [2]float64{x, y}, zoom, tilesize)

				centers = append(centers, [2]float64{lat, lng})

				r := image.Rect(
					int(x-(float64(tilesize/2)))-zeroPtX,
					int(y-(float64(tilesize/2)))-zeroPtY,
					int(x+(float64(tilesize/2)))-zeroPtX,
					int(y+(float64(tilesize/2)))-zeroPtY,
				)
				log.Printf("center: x %v y %v -- Lat: %v Lng: %v r: %v", x, y, lat, lng, r)
				rects = append(rects, r)
			}
			lat, lng := bounds.PointToLatLng(prj, [2]float64{x, y}, zoom, tilesize)
			centers = append(centers, [2]float64{lat, lng})

			r := image.Rect(
				int(x-(float64(tilesize/2)))-zeroPtX,
				int(y-(float64(tilesize/2)))-zeroPtY,
				int(x+(float64(tilesize/2)))-zeroPtX,
				int(y+(float64(tilesize/2)))-zeroPtY,
			)
			log.Printf("center: x %v y %v -- Lat: %v Lng: %v r: %v", x, y, lat, lng, r)
			rects = append(rects, r)

		}
		{
			x := startingPtX

			for ; x <= centerPtX+float64(cmdGenerateWidth/2); x += float64(tilesize) {
				lat, lng := bounds.PointToLatLng(prj, [2]float64{x, y}, zoom, tilesize)
				centers = append(centers, [2]float64{lat, lng})

				r := image.Rect(
					int(x-(float64(tilesize/2)))-zeroPtX,
					int(y-(float64(tilesize/2)))-zeroPtY,
					int(x+(float64(tilesize/2)))-zeroPtX,
					int(y+(float64(tilesize/2)))-zeroPtY,
				)
				log.Printf("center: x %v y %v -- Lat: %vLng: %v r: %v", x, y, lat, lng, r)
				rects = append(rects, r)
			}

			lat, lng := bounds.PointToLatLng(prj, [2]float64{x, y}, zoom, tilesize)
			centers = append(centers, [2]float64{lat, lng})

			r := image.Rect(
				int(x-(float64(tilesize/2)))-zeroPtX,
				int(y-(float64(tilesize/2)))-zeroPtY,
				int(x+(float64(tilesize/2)))-zeroPtX,
				int(y+(float64(tilesize/2)))-zeroPtY,
			)
			log.Printf("center: x %v y %v -- Lat: %v Lng: %v r: %v", x, y, lat, lng, r)
			rects = append(rects, r)
			maxWidth = int(x+(float64(tilesize/2))) - zeroPtX
			maxHeight = int(y+(float64(tilesize/2))) - zeroPtY
		}
	}

	log.Printf("Max Width %v, Height %v", maxWidth, maxHeight)
	dstimg := image.NewNRGBA(
		image.Rect(
			0,
			0,
			maxWidth,
			maxHeight,
		),
	)

	log.Printf("There are %v centers", len(centers))
	for i := range centers {
		snpsht := mbgl.Snapshotter{
			Style:    RootStyle,
			Width:    uint32(tilesize + padding),
			Height:   uint32(tilesize + padding),
			PPIRatio: 1.0,
			Lat:      centers[i][0],
			Lng:      centers[i][1],
			Zoom:     zoom,
		}
		img, err := mbgl.Snapshot1(snpsht)
		if err != nil {
			return nil, err
		}
		log.Printf("%v: Drawing tile: %v,%v %v w: %v h: %v", i, centers[i][0], centers[i][1], zoom, uint32(tilesize+10), uint32(tilesize+10))
		debugWriteoutimage(img, fmt.Sprintf("img_%v_c%v-%vz%vtz%v", i, centers[i][0], centers[i][1], zoom, tilesize+10))
		draw.Draw(dstimg, rects[i], img, image.ZP, draw.Over)
		debugWriteoutimage(dstimg, fmt.Sprintf("destimg_%v_c%v-%vz%vtz%v", i, centers[i][0], centers[i][1], zoom, tilesize+10))
	}

	return dstimg, nil
}
func debugWriteoutimage(img image.Image, n string) {
	name := n + ".jpg"

	w, err := os.Create(name)
	if err != nil {
		log.Printf("error creating output file: %v -- %v\n", name, err)
		return
	}
	if err := jpeg.Encode(w, img, nil); err != nil {
		log.Printf("error creating output file: %v -- %v\n", name, err)
		w.Close()
		return
	}
	w.Close()
}

func init() {
	pf := cmdGenerate.PersistentFlags()
	pf.IntVar(&cmdGenerateWidth, "width", 512, "Width of the image to generate.")
	pf.IntVar(&cmdGenerateHeight, "height", 512, "Height of the image to generate.")
	pf.Float64Var(&cmdGeneratePPIRatio, "ppiratio", 1.0, "The pixel per inch ratio.")
	pf.Float64Var(&cmdGeneratePitch, "pitch", 0.0, "The pitch of the map.")
	pf.Float64Var(&cmdGenerateBearing, "bearing", 0.0, "The bearing of the map.")
	pf.StringVar(&cmdGenerateFormat, "format", "", "Defaults to the ext of the output file, or jpg if not provided.")

	cmdGenerate.AddCommand(cmdGenerateBounds)
	cmdGenerate.AddCommand(cmdGenerateCenter)
}
