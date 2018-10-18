package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateCenter = &cobra.Command{
	Use:   "center lng lat zoom [output_file.jpg]",
	Short: "generate image using center and zoom",
	Long:  `use a center point (lng lat) and a zoom to generate the image.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(3, 4)(cmd, args); err != nil {
			return err
		}
		if _, ok := IsValidLngString(args[0]); !ok {
			return fmt.Errorf("Longitude must be between -180.0 and 180.0 : given %v", args[0])
		}
		if _, ok := IsValidLatString(args[1]); !ok {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[1])
		}
		z64, err := strconv.ParseFloat(strings.TrimSpace(args[2]), 64)
		if err != nil || z64 < 0 || z64 > 22 {
			return fmt.Errorf("zoom (%v) must be a number from 0 - 22", args[2])
		}

		if len(args) == 4 {
			cmdGenerateOutputName = args[3]
		}

		return nil
	},
	Run: commandGenerateCenterTiles,
}

func getFormatString(file string) (ext string, err error) {
	if cmdGenerateFormat != "" {
		ext = cmdGenerateFormat
	} else {
		ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		if ext == "" {
			ext = "jpg"
		} else if ext != "jpg" && ext != "png" {
			return "jpg", fmt.Errorf("output format(%v) must be jpg or png", ext)
		}
	}
	return ext, nil
}

func commandGenerateCenter(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)

	// Already checked in the Args validation function
	lng, _ := strconv.ParseFloat(strings.TrimSpace(args[0]), 64)
	lat, _ := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
	zoom, _ := strconv.ParseFloat(strings.TrimSpace(args[2]), 64)

	if err := genImage([2]float64{lng, lat}, zoom, cmdGenerateOutputName); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func commandGenerateCenterTiles(cmd *cobra.Command, args []string) {

	// 4096 is the tile size.
	const tilesize = 4096 / 2

	if math.Max(float64(cmdGenerateHeight), float64(cmdGenerateWidth)) <= tilesize {
		// we don't need to combine tiles, we can just return using the normal commandGenerateBounds
		commandGenerateCenter(cmd, args)
		return
	}

	// assume ESPG3857 for now.
	prj := bounds.ESPG3857

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)

	// Already checked in the Args validation function
	lng, _ := strconv.ParseFloat(strings.TrimSpace(args[0]), 64)
	lat, _ := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
	zoom, _ := strconv.ParseFloat(strings.TrimSpace(args[2]), 64)

	centerPt := bounds.LatLngToPoint(prj, lat, lng, zoom, tilesize)
	log.Printf("Going to draw: lng(%v),lat(%v) z%v -- %v x %v", lng, lat, zoom, cmdGenerateWidth, cmdGenerateHeight)

	err := getOutfile(cmdGenerateOutputName, func(out io.Writer, ext string) error {

		dstimg, err := genCenterZoomTilesImage(centerPt[0], centerPt[1], zoom, tilesize)
		if err != nil {
			return err
		}
		_, err = writeImage(out, dstimg, cmdGenerateWidth, cmdGenerateHeight, ext)
		return err

	})

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

}
