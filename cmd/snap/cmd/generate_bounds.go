package cmd

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/spherical"
	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateBoundsHull *geom.Extent

var cmdGenerateBounds = &cobra.Command{
	Use:   "bounds lng lat lng lat [output_file.jpg]",
	Short: "generate image using bounds",
	Long:  `use a bounds described as (lng lat lng lat) set of coordinates to generate the image.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(4, 5)(cmd, args); err != nil {
			return err
		}
		var lng1, lat1, lng2, lat2 float64
		var ok bool
		if lng1, ok = IsValidLngString(args[0]); !ok {
			return fmt.Errorf("Longitude must be between -180.0 and 180.0 : given %v", args[0])
		}
		if lat1, ok = IsValidLatString(args[1]); !ok {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[1])
		}
		if lng2, ok = IsValidLngString(args[2]); !ok {
			return fmt.Errorf("Longitude must be between -180.0 and 80.0 : given %v", args[2])
		}
		if lat2, ok = IsValidLatString(args[3]); !ok {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[3])
		}
		cmdGenerateBoundsHull = spherical.Hull([2]float64{lat1, lng1}, [2]float64{lat2, lng2})
		if len(args) == 5 {
			cmdGenerateOutputName = args[4]
		}
		return nil
	},
	Run: commandGenerateBoundsViaTiles,
}

func commandGenerateBounds(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)

	center, zoom := bounds.CenterZoom(cmdGenerateBoundsHull, float64(cmdGenerateWidth), float64(cmdGenerateHeight))

	if err := genImage(center, zoom, cmdGenerateOutputName); err != nil {
		cmd.Println(err.Error())
		os.Exit(1)
	}
}

func commandGenerateBoundsViaTiles(cmd *cobra.Command, args []string) {

	// 4096 is the tile size.
	const tilesize = 4096

	if math.Max(float64(cmdGenerateHeight), float64(cmdGenerateWidth)) <= tilesize {
		// we don't need to combine tiles, we can just return using the normal commandGenerateBounds
		commandGenerateBounds(cmd, args)
		return
	}

	// assume ESPG3857 for now.
	prj := bounds.ESPG3857
	hull := cmdGenerateBoundsHull
	if hull == nil {
		hull = prj.Bounds()
	}
	zoom := bounds.ZoomTile(hull, float64(cmdGenerateWidth), float64(cmdGenerateHeight), tilesize)

	// for lat lng geom.Extent should be laid out as follows:
	// {west, south, east, north}
	ne := [2]float64{hull[3], hull[2]}
	sw := [2]float64{hull[1], hull[0]}

	nePt := prj.Transform(prj.Project(ne), tilesize/256)
	swPt := prj.Transform(prj.Project(sw), tilesize/256)

	// center point.
	centerPtX := (swPt[0] + nePt[0]) / 2
	centerPtY := (swPt[1] + nePt[1]) / 2

	err := getOutfile(cmdGenerateOutputName, func(out io.Writer, ext string) error {

		dstimg, err := genCenterZoomTilesImage(centerPtX, centerPtY, zoom, tilesize)
		if err != nil {
			return err
		}
		_, err = writeImage(out, dstimg, cmdGenerateWidth, cmdGenerateHeight, ext)
		return err

	})

	if err != nil {
		cmd.Println(err.Error())
		os.Exit(1)
	}
}
