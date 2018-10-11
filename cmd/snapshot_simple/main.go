package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/geom/slippy"
	"github.com/go-spatial/geom/spherical"
	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
)

/*
snapshotter bounds "lat long lat long" -width=100, -height=100 style filename.png
snapshotter tile "z/x/y" style filename.png
*/

var FWidth uint
var FHeight uint
var FStyle string
var FPixelRatio float64
var FCenter [2]float64
var FZoom uint

var FOutputFilename string

func usage() {
	fmt.Fprintf(os.Stderr, "%v [options...] \"lat long lat long\" output_filename.png\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%v [options...] bounds \"lat long lat long\" output_filename.png\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%v [options...] tile \"z/x/y\" output_filename.png\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%v [options...] center \"lat long zoom\" output_filename.png\n", os.Args[0])
	flag.PrintDefaults()
}

func parseBounds(boundString string) {
	var err error
	bnds := strings.Split(boundString, " ")
	if len(bnds) != 4 {
		fmt.Fprintf(os.Stderr, "Error: invalid bounds provided %v\n", boundString)
		usage()
		os.Exit(2)
	}
	var fbounds [4]float64
	for i, bound := range bnds {
		fbounds[i], err = strconv.ParseFloat(strings.TrimSpace(bound), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid bounds provided %v\n", boundString)
			fmt.Fprintf(os.Stderr, "Error: unabled to parse %v(%v) as a float.\n", bound, i)
			usage()
			os.Exit(2)
		}
	}
	// our fbounds is in lng lat lng lat order need to fix. to lat lng lat lng
	hull := spherical.Hull([2]float64{fbounds[1], fbounds[0]}, [2]float64{fbounds[3], fbounds[2]})
	var zoom float64
	FCenter, zoom = bounds.CenterZoom(hull, float64(FWidth), float64(FHeight))
	FZoom = uint(zoom)
}

func parseTile(tileString string) {
	var err error
	var v uint64
	parts := strings.Split(tileString, "/")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Error: invalid z/x/y coordinates %v\n", tileString)
		usage()
		os.Exit(2)
	}
	var label = [...]string{"Z", "X", "Y"}
	var fTile slippy.Tile

	for i, part := range parts {
		v, err = strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid tile coordinates provided %v\n", tileString)
			fmt.Fprintf(os.Stderr, "Error: unabled to parse %v %v as a uint.\n", label[i], part)
			usage()
			os.Exit(2)
		}
		switch i {
		case 0:
			fTile.Z = uint(v)
		case 1:
			fTile.X = uint(v)
		case 2:
			fTile.Y = uint(v)
		}
	}
	FCenter = bounds.Center(fTile.Extent3857())
	FZoom = fTile.Z
}

func parseCenterZoom(centerZoomString string) {
	var err error
	errFn := func(label string, item string, expectedType string, err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid center zoom provided %v\n", centerZoomString)
			fmt.Fprintf(os.Stderr, "Error: unabled to parse %v(%v) as a %v.\n", label, expectedType, item)
			usage()
			os.Exit(2)
		}
	}
	centerzoom := strings.Split(centerZoomString, " ")
	if len(centerzoom) != 3 {
		fmt.Fprintf(os.Stderr, "Error: invalid center zoom provided %v\n", centerZoomString)
		usage()
		os.Exit(2)
	}

	FCenter[0], err = strconv.ParseFloat(strings.TrimSpace(centerzoom[0]), 64)
	errFn("lat", centerzoom[0], "float", err)
	FCenter[1], err = strconv.ParseFloat(strings.TrimSpace(centerzoom[1]), 64)
	errFn("lng", centerzoom[1], "float", err)
	zoom, err := strconv.ParseUint(strings.TrimSpace(centerzoom[2]), 10, 64)
	errFn("zoom", centerzoom[2], "uint", err)
	FZoom = uint(zoom)

}

func ParseFlags() {

	flag.UintVar(&FWidth, "width", 512, "Width of the image to generate.")
	flag.UintVar(&FWidth, "w", 512, "Width of the image to generate.")

	flag.UintVar(&FHeight, "height", 512, "Height of the image to generate.")
	flag.UintVar(&FHeight, "h", 512, "Height of the image to generate.")

	flag.StringVar(&FStyle, "style", "file://style.json", "Style file")
	flag.Float64Var(&FPixelRatio, "pixel", 1.0, "The pixel ratio")

	flag.Parse()

	if flag.NArg() < 3 {
		usage()
		os.Exit(2)
	}

	var fileIdx = 2

	switch strings.TrimSpace(strings.ToLower(flag.Arg(0))) {
	case "bounds":
		// The next variable should be the bounds seperated by spaces.
		parseBounds(flag.Arg(1))
	case "tile":
		// The next variable should be the coordinates sepearted by forward slash.
		parseTile(flag.Arg(1))
	case "center":
		parseCenterZoom(flag.Arg(1))
	default: // assume bounds as the default subcommand
		parseBounds(flag.Arg(0))
		fileIdx = 1
	}

	// Next should be the output filename.
	FOutputFilename = strings.TrimSpace(flag.Arg(fileIdx))
}

func main() {
	var file *os.File
	var err error
	var img image.Image
	ParseFlags()
	mbgl.NewRunLoop()
	defer mbgl.DestroyRunLoop()
	snpsht := mbgl.Snapshotter{
		Style:    FStyle,
		Width:    uint32(FWidth),
		Height:   uint32(FHeight),
		PPIRatio: FPixelRatio,
		Lat:      FCenter[0],
		Lng:      FCenter[1],
		Zoom:     float64(FZoom),
	}
	img, err = mbgl.Snapshot(snpsht)
	if err != nil {
		fmt.Fprintf(os.Stderr, "got an error creating the snapshot: %v\n", err)
		os.Exit(3)
	}
	file, err = os.Create(FOutputFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output file: %v -- %v\n", FOutputFilename, err)
		os.Exit(3)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		log.Printf("Failed to write %v", FOutputFilename)
		log.Println(err)
		return
	}

	fmt.Printf("successfully wrote outfile: %v\n", FOutputFilename)
}
