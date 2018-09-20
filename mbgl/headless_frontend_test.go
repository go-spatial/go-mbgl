package mbgl

import (
	"fmt"
	"runtime"
	"testing"
	"os"

	"github.com/arolek/p"
)

func TestNewHeadlessFrontend(t *testing.T) {
	type tcase struct {
		size       Size
		pixelRatio float32
		src        FileSource
	}

	fn := func(tc tcase, t *testing.T) {
		runtime.LockOSThread()
		fmt.Fprintf(os.Stderr,"Starting new run loop\n")
		loop := NewRunLoop()
		fmt.Fprintf(os.Stderr,"Setting up new thread pool 4\n")
		tpool := NewThreadPool(4)


		fmt.Fprintf(os.Stderr,"Setting new Headless frontend\n")
		frontend := NewHeadlessFrontend(&tc.size, tc.pixelRatio, tc.src, tpool, nil, nil)

		fmt.Fprintf(os.Stderr,"Setting new Headless frontend: %p\n",frontend)

		fmt.Fprintf(os.Stderr,"Destroying new Headless frontend\n")

		frontend.Destruct()

		fmt.Fprintf(os.Stderr,"Destroy the testcase size.\n")
		tc.size.Destruct()
		fmt.Fprintf(os.Stderr,"Destroy the testcase source.\n")
		tc.src.Destruct()

		fmt.Fprintf(os.Stderr,"Destroy the pool\n")
		tpool.Destruct()

		fmt.Fprintf(os.Stderr,"Destroy the loop\n")
		loop.Destruct()

		fmt.Fprintf(os.Stderr,"Unlock the thread\n")
		runtime.UnlockOSThread()
	}

	testcases := map[string]tcase{
		"1": {
			size:       Size{Height: 256, Width: 256},
			pixelRatio: 1.0,
			src: NewDefaultFileSource("",
				"https://osm.tegola.io/maps/osm/style.json",
				p.Uint64(0)),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}

func TestHeadlessFrontendRender(t *testing.T) {
	type tcase struct {
		size       Size
		pixelRatio float32
		src        FileSource
	}

	fn := func(tc tcase, t *testing.T) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		loop := NewRunLoop()
		defer loop.Destruct()

		tpool := NewThreadPool(4)
		defer tpool.Destruct()

		hfe := NewHeadlessFrontend(&tc.size,
			tc.pixelRatio,
			tc.src,
			tpool,
			nil, nil)
		defer hfe.Destruct()

		m := NewMap(hfe,
			tc.size,
			tc.pixelRatio,
			tc.src,
			tpool)
		defer m.Destruct()

		m.setStyleUrl("https://osm.tegola.io/maps/osm/style.json")

		fmt.Println("calling render")
		hfe.Render(m)
		fmt.Println("rendered")

		tc.size.Destruct()
		tc.src.Destruct()
	}

	testcases := map[string]tcase{
		"1": {
			size:       Size{Height: 256, Width: 256},
			pixelRatio: 1.0,
			src:        NewDefaultFileSource("", "", p.Uint64(0)),
		},
	}

	for k, v := range testcases {
		t.Run(k, func(t *testing.T) {
			fn(v, t)
		})
	}
}
