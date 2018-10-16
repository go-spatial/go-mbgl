package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-spatial/geom/spherical"
	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateBounds = &cobra.Command{
	Use:   "bounds lng lat lng lat [output_file.jpg]",
	Short: "generate image using bounds",
	Long:  `use a bounds described as (lng lat lng lat) set of coordinates to generate the image.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 4 {
			return errors.New("requres at bounds( lng lat lng lat )")
		}
		if len(args) > 5 {
			return errors.New("extra values provided")
		}
		if !IsValidLngString(args[0]) {
			return fmt.Errorf("Longitude must be between -180.0 and 180.0 : given %v", args[0])
		}
		if !IsValidLatString(args[1]) {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[1])
		}
		if !IsValidLngString(args[2]) {
			return fmt.Errorf("Longitude must be between -180.0 and 80.0 : given %v", args[2])
		}
		if !IsValidLatString(args[3]) {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[3])
		}
		return nil
	},
	Run: commandGenerateBounds,
}

func commandGenerateBounds(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)

	// Already checked in the Args validation function
	lng1, _ := strconv.ParseFloat(strings.TrimSpace(args[0]), 64)
	lat1, _ := strconv.ParseFloat(strings.TrimSpace(args[1]), 64)
	lng2, _ := strconv.ParseFloat(strings.TrimSpace(args[2]), 64)
	lat2, _ := strconv.ParseFloat(strings.TrimSpace(args[3]), 64)
	hull := spherical.Hull([2]float64{lat1, lng1}, [2]float64{lat2, lng2})
	center, zoom := bounds.CenterZoom(hull, float64(cmdGenerateWidth), float64(cmdGenerateHeight))
	output := ""
	if len(args) == 5 {
		output = args[4]
	}
	if err := genImage(center, zoom, output); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
