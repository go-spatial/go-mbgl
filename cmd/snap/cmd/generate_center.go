package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"
	"github.com/spf13/cobra"
)

var cmdGenerateCenter = &cobra.Command{
	Use:   "center lng lat zoom [output_file.jpg]",
	Short: "generate image using center and zoom",
	Long:  `use a center point (lng lat) and a zoom to generate the image.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("requres a center (lng lat) and a zoom.")
		}
		if len(args) > 4 {
			return errors.New("extra values provided")
		}
		if !IsValidLngString(args[0]) {
			return fmt.Errorf("Longitude must be between -180.0 and 180.0 : given %v", args[0])
		}
		if !IsValidLatString(args[1]) {
			return fmt.Errorf("Latitude must be between -90.0 and 90.0 : given %v", args[1])
		}
		z64, err := strconv.ParseFloat(strings.TrimSpace(args[2]), 64)
		if err != nil || z64 < 0 || z64 > 22 {
			return fmt.Errorf("zoom (%v) must be a number from 0 - 22", args[2])
		}

		return nil
	},
	Run: commandGenerateCenter,
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

	output := ""
	if len(args) == 4 {
		output = args[3]
	}
	if err := genImage([2]float64{lng, lat}, zoom, output); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
