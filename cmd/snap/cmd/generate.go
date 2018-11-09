package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/spherical"
	"github.com/go-spatial/go-mbgl/cmd/snap/generate"
	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var (
	// cmdGenerateWidth is the desired width of the image, this is not the actual width which will be based
	// on this value and the PPI ratio.
	cmdGenerateWidth int
	// cmdGenerateHeight is the desired height of the image, this is not the actual height which will be based
	// on this value and the PPI ratio.
	cmdGenerateHeight int
	// cmdGeneatePPIRatio is the Pixel Per inch ratio. Default is 1.0
	cmdGeneratePPIRatio float64
	// cmdGeneratePitch is the pitch of the map when being rendered.
	cmdGeneratePitch float64
	// cmdGenerateBearing is the bearing of the map when being rendered.
	cmdGenerateBearing float64
	// cmdGenerateFormat the output format to use.
	cmdGenerateFormat string
	// cmdGenerateOutputName is the name of the file to generate.
	cmdGenerateOutputName string
	// cmdGenerateCenter is the parameter for specifying a center/zoom combo  (lng,lat,z)
	cmdGenerateCenter string
	// cmdGenerateBounds is the parameter for specifying a bounds to use to generate an image (lng,lat,lng,lat)
	cmdGenerateBounds     string
	cmdGenerateCenterZoom [3]float64 // lat, lng, z

	// zoom for bounds command.
	cmdGenerateZoom float64
)

var debugHull *geom.Extent

var cmdGenerate = &cobra.Command{
	Use:     "generate",
	Short:   "generate an image for given map coordinates.",
	Aliases: []string{"gen"},
	Args:    ValidateGenerateParams,
	Run:     commandGenerateCenterZoom,
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

func ValidateGenerateParams(cmd *cobra.Command, args []string) (err error) {
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
	if cmdGenerateCenter == "" && cmdGenerateBounds == "" {
		return errors.New("--center or --bounds must be specified.")
	}
	if cmdGenerateCenter != "" && cmdGenerateBounds != "" {
		return errors.New("only one --center or --bounds must be specified.")
	}
	if cmdGenerateCenter != "" {
		var ok bool
		// parse the center and zoom, which should be seperated by commas
		parts := strings.Split(cmdGenerateCenter, ",")
		if len(parts) < 3 {
			return errors.New("need lng,lat and zoom values for center.")
		}
		if cmdGenerateCenterZoom[0], ok = IsValidLatString(parts[1]); !ok {
			return fmt.Errorf("invalid lat expected value to be between -90 and 90: %v", parts[1])
		}
		if cmdGenerateCenterZoom[1], ok = IsValidLngString(parts[0]); !ok {
			return fmt.Errorf("invalid lng expected value to be between -180 and 180: %v", parts[0])
		}
		if cmdGenerateCenterZoom[2], err = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64); err != nil {
			return fmt.Errorf("invalid zoom: %v", err)
		}
	}
	if cmdGenerateBounds != "" {
		// lng,lat,lng,lat
		// parse the lng,lat,lng,lat, which should be seperated by commas
		parts := strings.Split(cmdGenerateBounds, ",")
		var coords [2][2]float64
		var ok bool
		if len(parts) < 4 {
			return errors.New("need lng,lat,lng,lat for bounds.")
		}
		if coords[0][0], ok = IsValidLngString(parts[0]); !ok {
			return fmt.Errorf("invalid lng expected value to be between -90 and 90: %v", parts[0])
		}
		if coords[0][1], ok = IsValidLatString(parts[1]); !ok {
			return fmt.Errorf("invalid lat expected value to be between -180 and 180: %v", parts[1])
		}
		if coords[1][0], ok = IsValidLngString(parts[2]); !ok {
			return fmt.Errorf("invalid lng expected value to be between -90 and 90: %v", parts[2])
		}
		if coords[1][1], ok = IsValidLatString(parts[3]); !ok {
			return fmt.Errorf("invalid lat expected value to be between -180 and 180: %v", parts[3])
		}
		hull := spherical.Hull(coords[0], coords[1])
		var center [2]float64

		// Next we need to see if we have a zoom or width and height. We do this by checking  if zoom
		// was explicitly set. If so, then we need to make sure width and height weren't set.
		zset := cmd.Flag("zoom").Changed
		whset := cmd.Flag("width").Changed || cmd.Flag("height").Changed
		if zset && whset {
			return fmt.Errorf("can not specify zoom along with width or height parameters.")
		}

		if zset {
			cmdGenerateCenterZoom[2] = cmdGenerateZoom
			center = bounds.Center(hull, cmdGenerateZoom)
			// need to set the width and height based on the bounds.
			width, height := bounds.WidthHeightTile(hull, cmdGenerateZoom, 4096/8)
			// 4 is the scale factor
			cmdGenerateWidth = int(width)
			cmdGenerateHeight = int(height)

			debugHull = hull

		} else {
			center, cmdGenerateCenterZoom[2] = bounds.CenterZoom(hull, float64(cmdGenerateWidth), float64(cmdGenerateHeight))
		}

		cmdGenerateCenterZoom[0], cmdGenerateCenterZoom[1] = center[0], center[1]
	}

	if len(args) == 1 {
		cmdGenerateOutputName = args[0]
	}

	return nil
}

func commandGenerateCenterZoom(cmd *cobra.Command, args []string) {

	const tilesize = 4096 / 2

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mbgl.StartSnapshotManager(ctx)

	err := getOutfile(cmdGenerateOutputName, func(out io.Writer, ext string) error {
		prj := bounds.ESPG3857

		centerPt := bounds.LatLngToPoint(prj, cmdGenerateCenterZoom[0], cmdGenerateCenterZoom[1], cmdGenerateCenterZoom[2], tilesize)
		dstimg, err := generate.NewImage(
			prj,
			cmdGenerateWidth, cmdGenerateHeight,
			centerPt,
			cmdGenerateCenterZoom[2],
			cmdGeneratePPIRatio,
			cmdGeneratePitch,
			cmdGenerateBearing,
			RootStyle,
			"", "",
		)
		if err != nil {
			return err
		}
		if debugHull != nil {
			dstimg.SetDebugBounds(debugHull, cmdGenerateCenterZoom[2])
		}
		defer dstimg.Close()
		_, err = writeImage(out, dstimg, cmdGenerateWidth, cmdGenerateHeight, ext)
		return err

	})

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

}

func IsValidLngString(lng string) (float64, bool) {

	f64, err := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if err != nil {
		log.Println("Got error Parsing ", err)
		return f64, false
	}
	return f64, -180.0 <= f64 && f64 <= 180.0
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

func genCenterZoomTilesImage(centerPtX, centerPtY, zoom, ppiratio, pitch, bearing float64, style string) (*generate.Image, error) {

	return generate.NewImage(
		bounds.ESPG3857,
		cmdGenerateWidth, cmdGenerateHeight,
		[2]float64{centerPtX, centerPtY}, zoom,
		ppiratio,
		pitch,
		bearing,
		style,
		"", "",
	)
}

func init() {

	/*

	 usage

	 snap --bounds --width --height --ppiratio

	 snap --bounds --zoom

	 snap --bounds

	 invalid

	 snap --bounds --zoom --width --height


	*/

	pf := cmdGenerate.PersistentFlags()
	pf.IntVar(&cmdGenerateWidth, "width", 512, "Width of the image to generate.")
	pf.IntVar(&cmdGenerateHeight, "height", 512, "Height of the image to generate.")
	pf.Float64Var(&cmdGenerateZoom, "zoom", 0.0, "zoom value")

	pf.Float64Var(&cmdGeneratePPIRatio, "ppiratio", 1.0, "The pixel per inch ratio.")
	pf.Float64Var(&cmdGeneratePitch, "pitch", 0.0, "The pitch of the map.")

	pf.Float64Var(&cmdGenerateBearing, "bearing", 0.0, "The bearing of the map.")
	pf.StringVar(&cmdGenerateFormat, "format", "", "Defaults to the ext of the output file, or jpg if not provided.")
	pf.StringVar(&cmdGenerateCenter, "center", "", "Generate the image based on the center: lat,lng,z ")
	pf.StringVar(&cmdGenerateBounds, "bounds", "", "Generate the image based on the bounds: nelng,nelat,swlng,swlat")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
