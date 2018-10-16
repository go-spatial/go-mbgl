package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var cmdGenerateWidth int
var cmdGenerateHeight int
var cmdGeneratePPIRatio float64
var cmdGeneratePitch float64
var cmdGenerateBearing float64
var cmdGenerateFormat string

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

func IsValidLngString(lng string) bool {

	f64, err := strconv.ParseFloat(strings.TrimSpace(lng), 64)
	if err != nil {
		return false
	}
	return -180.0 <= f64 && f64 <= 190.0
}

func IsValidLatString(lat string) bool {

	f64, err := strconv.ParseFloat(strings.TrimSpace(lat), 64)
	if err != nil {
		return false
	}
	return -90.0 <= f64 && f64 <= 90.0
}

func genImage(center [2]float64, zoom float64, output string) error {
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

	if _, err := generateZoomCenterImage(out, RootStyle, cmdGenerateWidth, cmdGenerateHeight, cmdGeneratePPIRatio, cmdGeneratePitch, cmdGenerateBearing, center, zoom, ext); err != nil {
		return fmt.Errorf("failed to generate image: %v", err)
	}
	return nil
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
