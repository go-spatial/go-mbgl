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

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/slippy"
	mbgl "github.com/go-spatial/go-mbgl"
)

/*
snapshotter bounds "lat long lat long" -width=100, -height=100 style filename.png
snapshotter tile "z/x/y" style filename.png
*/

type cmdType uint8

const (
	CmdUnknown = cmdType(iota)
	CmdBounds
	CmdTile
)

func (c cmdType) String() string {
	switch c {
	case CmdBounds:
		return "Bounds"
	case CmdTile:
		return "Tile"
	default:
		return "Unknown"
	}
}

var FWidth uint
var FHeight uint
var FStyle string
var FPixelRatio float64

var FBounds geom.Extent
var FTile slippy.Tile
var FOutputFilename string

func usage() {
	fmt.Fprintf(os.Stderr, "%v [options...] \"lat long lat long\" output_filename.png\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%v [options...] bounds \"lat long lat long\" output_filename.png\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "%v [options...] tile \"z/x/y\" output_filename.png\n", os.Args[0])
	flag.PrintDefaults()
}

func parseBounds(boundString string) {
	var err error
	bounds := strings.Split(boundString, " ")
	if len(bounds) != 4 {
		fmt.Fprintf(os.Stderr, "Error: invalid bounds provided — %v\n", boundString)
		usage()
		os.Exit(2)
	}
	for i, bound := range bounds {
		FBounds[i], err = strconv.ParseFloat(strings.TrimSpace(bound), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid bounds provided — %v\n", boundString)
			fmt.Fprintf(os.Stderr, "Error: unabled to parse %v(%v) as a float.\n", bound, i)
			usage()
			os.Exit(2)
		}
	}
}

func parseTile(tileString string) {
	var err error
	var v uint64
	parts := strings.Split(tileString, "/")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Error: invalid z/x/y coordinates  — %v\n", tileString)
		usage()
		os.Exit(2)
	}
	var label = [...]string{"Z", "X", "Y"}

	for i, part := range parts {
		v, err = strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid tile coordinates provided — %v\n", tileString)
			fmt.Fprintf(os.Stderr, "Error: unabled to parse %v %v as a uint.\n", label[i], part)
			usage()
			os.Exit(2)
		}
		switch i {
		case 0:
			FTile.Z = uint(v)
		case 1:
			FTile.X = uint(v)
		case 2:
			FTile.Y = uint(v)
		}
	}
}

func ParseFlags() cmdType {

	flag.UintVar(&FWidth, "width", 100, "Width of the image to generate.")
	flag.UintVar(&FWidth, "w", 100, "Width of the image to generate.")

	flag.UintVar(&FHeight, "height", 100, "Height of the image to generate.")
	flag.UintVar(&FHeight, "h", 100, "Height of the image to generate.")

	flag.StringVar(&FStyle, "style", "https://osm.tegola.io/capabilities/osm.json", "Style file")
	flag.Float64Var(&FPixelRatio, "pixel", 1.0, "The pixel ratio")

	flag.Parse()

	if flag.NArg() < 3 {
		usage()
		os.Exit(2)
	}

	var cmd cmdType
	var fileIdx = 2

	switch strings.TrimSpace(strings.ToLower(flag.Arg(0))) {
	case "bounds":
		// The next variable should be the bounds seperated by spaces.
		parseBounds(flag.Arg(1))
		cmd = CmdBounds
	case "tile":
		// The next variable should be the coordinates sepearted by forward slash.
		parseTile(flag.Arg(1))

		cmd = CmdTile
	default: // assume bounds as the default subcommand
		parseBounds(flag.Arg(0))
		fileIdx = 1
		cmd = CmdBounds
	}

	// Next should be the output filename.
	FOutputFilename = strings.TrimSpace(flag.Arg(fileIdx))
	return cmd

}

func main() {
	var file *os.File
	var err error
	var img image.Image
	cmd := ParseFlags()
	file, err = os.Create(FOutputFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
		os.Exit(3)
	}
	ss := mbgl.NewSnapshotter(FStyle, float32(FPixelRatio))
	size := image.Point{
		X: int(FWidth),
		Y: int(FHeight),
	}
	switch cmd {
	case CmdBounds:
		img = ss.Snapshot(&FBounds, size)

	case CmdTile:
		img = mbgl.SnapshotTile(ss, FTile, size)

	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		log.Println("Failed to write %v", FOutputFilename)
		log.Fatal(err)
	}
	file.Close()
	fmt.Printf("successfully wrote outfile: %v\n", FOutputFilename)
}
