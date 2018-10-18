package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/geom/slippy"
	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateTileTile slippy.Tile

var cmdGenerateTile = &cobra.Command{
	Use:   "tile z/x/y [output_file.jpg]",
	Short: "generate image using tile coordinates.",
	Long:  `use a tile coordinate (z/x/y) to generate the image.`,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			return err
		}
		z, x, y, err := parseTileCoordinate(args[0])
		if err != nil {
			return err
		}
		cmdGenerateTileTile.Z = uint(z)
		cmdGenerateTileTile.X = uint(x)
		cmdGenerateTileTile.Y = uint(y)
		return nil
	},
	Run: commandGenerateTile,
}

// parseTileCoordinate take a tile described as z/x/y and returns the compontents if it was able to parse the
// the string, or an error. Zoom
func parseTileCoordinate(crd string) (z, x, y int, err error) {
	coords := strings.Split(crd, "/")
	if len(coords) < 3 {
		return 0, 0, 0, fmt.Errorf("expected at least three components: %v", crd)
	}
	z64, err := strconv.ParseUint(strings.TrimSpace(coords[0]), 64, 10)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unable to parse z value as an int: %v", err)
	}
	x64, err := strconv.ParseInt(strings.TrimSpace(coords[1]), 64, 10)
	if err != nil {
		return int(z64), 0, 0, fmt.Errorf("unable to parse x value as an int: %v", err)
	}
	y64, err := strconv.ParseInt(strings.TrimSpace(coords[2]), 64, 10)
	if err != nil {
		return int(z64), int(x64), 0, fmt.Errorf("unable to parse y value as an int: %v", err)
	}
	return int(z64), int(x64), int(y64), nil
}

func commandGenerateTile(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)
	center := bounds.Center(cmdGenerateTileTile.Extent3857())
	zoom := cmdGenerateTileTile.Z
	output := ""
	if len(args) == 2 {
		output = args[1]
	}
	if err := genImage(center, float64(zoom), output); err != nil {
		cmd.Println(err.Error())
		os.Exit(1)
	}

}
