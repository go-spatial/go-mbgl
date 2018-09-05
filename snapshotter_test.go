package mbgl

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/slippy"
)

func tfile(name string) (*os.File, error) {
	fname := os.DevNull
	evar := os.Getenv("MBGL_TEST_OUT_DIR")
	if evar != "" {
		fmt.Println("outputing to ", evar)
		os.MkdirAll(evar, 0600)

		fname = strings.Replace(name, "/", "-", -1)
		fname = filepath.Join(evar, fname+".png")
	}

	return os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func TestSnapshot(t *testing.T) {
	type tcase struct {
		src   string
		ext   *geom.Extent
		size  image.Point
		ratio float32
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			img := Snapshot(tc.src, tc.ext, tc.size, tc.ratio)

			f, err := tfile(t.Name())
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			defer f.Close()

			if err = png.Encode(f, img); err != nil {
				t.Fatalf("unexpected error %v", err)
			}
		}
	}

	exts := make(map[string]*geom.Extent)

	tile := slippy.NewTileLatLon(12, 33, 117)

	exts["1"] = tile.Extent4326()

	z, x, y := tile.ZXY()
	tile = slippy.NewTile(z, x+1, y)
	exts["2"] = tile.Extent4326()

	fmt.Println(exts["1"])

	testcases := map[string]tcase{
		"1": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(1000, 1000),
			ratio: 1.0,
		},
		"2": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(1000, 1000),
			ratio: 2.0,
		},
		"3": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(1000, 1000),
			ratio: 0.5,
		},
		"4": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(500, 500),
			ratio: 2.0,
		},
		"5": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(500, 500),
			ratio: 0.5,
		},
		"6": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(500, 500),
			ratio: 1.0,
		},
		"7": {
			src:   "https://osm.tegola.io/maps/osm/style.json",
			ext:   exts["1"],
			size:  image.Pt(2000, 2000),
			ratio: 1.0,
		},
		"small": {
			src:  "https://osm.tegola.io/maps/osm/style.json",
			ext:  exts["1"],
			size: image.Pt(1, 1),
		},
	}

	for name, tc := range testcases {
		t.Run(name, fn(tc))
	}
}

func TestSnapshotTile(t *testing.T) {
	type tcase struct {
		src  string
		tile *slippy.Tile
		size image.Point
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			s := NewSnapshotter(tc.src, 1.0)
			img := SnapshotTile(s, *tc.tile, tc.size)

			f, err := tfile(t.Name())
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			defer f.Close()
			if png.Encode(f, img); err != nil {
				t.Fatalf("unxepected error %v", err)
			}
		}
	}

	testcases := map[string]tcase{
		"1": {
			src:  "https://osm.tegola.io/maps/osm/style.json",
			tile: slippy.NewTile(5, 9, 12),
			size: image.Pt(255, 255),
		},
		"2": {
			src:  "https://osm.tegola.io/maps/osm/style.json",
			tile: slippy.NewTile(5, 9, 12),
			size: image.Pt(512, 512),
		},
	}

	for name, tc := range testcases {
		t.Run(name, fn(tc))
	}
}
