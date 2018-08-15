package mbgl

import (
	"testing"
	"image/png"
	"os"
	"fmt"
	"github.com/go-spatial/geom/slippy"
	"github.com/go-spatial/geom"
)

func TestNewMapSnapshotter(t *testing.T) {
	type tcase struct {
		src        FileSource
		sched      Scheduler
		style      string
		size       Size
		pixelRatio float32
		camOpts    *CameraOptions
		region     *LatLngBounds
		cacheDir   *string
	}

	fn := func(tc tcase, t *testing.T) {
		ms := NewMapSnapshotter(
			tc.src,
			tc.sched,
			tc.style,
			tc.size,
			tc.pixelRatio,
			tc.camOpts,
			tc.region,
			tc.cacheDir)

		ms.Destruct()
	}

	testcases := map[string]tcase{
		"1": {
			src:        NewDefaultFileSource("", "", nil),
			sched:      NewThreadPool(4),
			style:      "https://osm.tegola.io/maps/osm/style.json",
			size:       Size{Width: 100, Height: 100},
			pixelRatio: 1.0,
			camOpts:    nil,
			region:     nil,
			cacheDir:   nil,
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}

func TestSnapshotterSnapshot(t *testing.T) {
	type tcase struct {
		ms *MapSnapshotter
	}

	fn := func(tc tcase, t *testing.T) {
		cImg := tc.ms.Snapshot()
		img := cImg.Image()

		fname := os.DevNull
		if evar := os.Getenv("MBGL_TEST_ON_DISK"); evar != "" {
			fmt.Println("outputing to ",evar)
			fname = evar
		}

		f, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0600)
		if err != nil {
			t.Fatalf("unexpected errro %v", err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			t.Fatalf("unexpected errro %v", err)
		}

		tc.ms.Destruct()
	}

	tpool := NewThreadPool(4)
	SchedulerSetCurrent(tpool)

	testcases := map[string]tcase{
		"1": {
			ms: NewMapSnapshotter(
				NewDefaultFileSource("", "", nil),
				tpool,
				"https://osm.tegola.io/maps/osm/style.json",
				Size{Height: 100, Width: 100},
				1.0,
				nil,
				nil,
				nil),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}

func TestSnapshotterSetCamOpts(t *testing.T) {
	type tcase struct {
		opts []CameraOptions
	}

	fn := func (tc tcase, t *testing.T) {
		tpool := NewThreadPool(4)
		SchedulerSetCurrent(tpool)

		ms := NewMapSnapshotter(
			NewDefaultFileSource("", "", nil),
			tpool,
			"https://osm.tegola.io/maps/osm/style.json",
			Size{Height: 100, Width: 100},
			1.0,
			&tc.opts[0],
			nil,
			nil)

		defer ms.Destruct()

		ms.Snapshot()

		for _, v := range tc.opts[1:] {
			ms.SetCameraOptions(v)
			ms.Snapshot()
		}
	}

	testcases := map[string]tcase {
		"1": {
			opts: []CameraOptions{
				{},
				{
					Center: NewLatLng(33, 117),
					Padding: NewEdgeInsets(10, 10, 10, 10),
				},
			},
		},
	}

	for k, v := range testcases {
		t.Run(k, func (t *testing.T) {
			fn(v, t)
		})
	}
}


func TestSnapshotterSetRegion(t *testing.T) {
	type tcase struct {
		bounds []*LatLngBounds
	}

	fn := func(tc tcase, t * testing.T) {
		tpool := NewThreadPool(4)
		SchedulerSetCurrent(tpool)

		ms := NewMapSnapshotter(
			NewDefaultFileSource("","", nil),
			tpool,
			"https://osm.tegola.io/maps/osm/style.json",
			Size{Height: 100, Width: 100},
			1.0,
			nil,
			tc.bounds[0],
			nil)
		defer ms.Destruct()

		ms.Snapshot()


		for _, v := range tc.bounds[1:] {
			ms.SetRegion(v)
			ms.Snapshot()
		}
	}

	testcases := map[string]tcase {
		"1" : {
			bounds: []*LatLngBounds {
				NewLatLngBoundsFromTile(slippy.NewTile(0, 0, 0, 0, geom.WebMercator)),
				NewLatLngBoundsFromTile(slippy.NewTile(12, 212, 6079, 0, geom.WebMercator)),
			},
		},
		"2": {
			bounds: []*LatLngBounds {
				nil,
				NewLatLngBounds(NewLatLng(33, 117), NewLatLng(34, 118)),
			},
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}