package mbgl

import (
	"image"
	"runtime"
	"sync"

	"github.com/arolek/p"
	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/slippy"
	"github.com/go-spatial/go-mbgl/mbgl"
)

func init() {
	runtime.GOMAXPROCS(1)
}

// This will set the size of the mbgl::ThreadPool
var DefaultThreadPoolSize = 4

// Snapshotter is an interface for processing raster image from
// vector tiles.
type Snapshotter interface {
	Snapshot(style string, extent *geom.Extent, size image.Point) image.Image
}

type snapshotterPool struct {
	pool sync.Pool
}

func (sp snapshotterPool) Snapshot(style string, extent *geom.Extent, size image.Point) image.Image {
	v := sp.pool.Get()
	vv := v.(Snapshotter)

	img := vv.Snapshot(style, extent, size)

	sp.pool.Put(v)
	return img
}

func NewSnapshotterPool(pixelRatio float32) Snapshotter {

	return snapshotterPool{
		pool: sync.Pool{
			New: func() interface{} {
				return NewSnapshotter(pixelRatio)
			},
		},
	}
}

// this is an internal representation o
type snapshotter struct {
	fsrc mbgl.FileSource // an interface
	snap *mbgl.MapSnapshotter
	size mbgl.Size

	// the mbgl::MapSnapshotter.Snapshot method is not thread safe
	snapLock sync.Mutex
}

func (s snapshotter) destruct() {
	s.fsrc.Destruct()
	s.snap.Destruct()
	s.size.Destruct()
}

// High level function for taking a single snapshot. This simply creates a new snapshotter and uses it.
func Snapshot(style string, ext *geom.Extent, size image.Point, pixelRatio float32) image.Image {
	return NewSnapshotter(pixelRatio).Snapshot(style, ext, size)
}

// This creates an instance of a Snapshotter with the specified style.
// Note: this high level implementation is thread safe, but performance might be better to lower the DefaultThreadPoolSize and use multiple snapshotters
// TODO(@ear7h): write benchmarks
func NewSnapshotter(pixelRatio float32) Snapshotter {
	src := mbgl.NewDefaultFileSource("", "", p.Uint64(0))

	tpool := mbgl.NewThreadPool(DefaultThreadPoolSize)
	mbgl.SchedulerSetCurrent(tpool)

	size := mbgl.Size{Width: 100, Height: 100}

	if pixelRatio == 0 {
		pixelRatio = 1.0
	}

	snap := mbgl.NewMapSnapshotter(src,
		tpool,
		"",
		size,
		pixelRatio,
		nil,
		nil,
		nil)

	ret := &snapshotter{
		fsrc: src,
		snap: snap,
		size: size,
	}

	// finalizer has to be on Go composite objecte because mbgl
	// types are pointers to empty structs as far as go knows
	runtime.SetFinalizer(ret, (*snapshotter).destruct)

	return ret
}

// Take the snapshot with a *geom.Extent encoded as lat/lng (ie. WSG84).
func (s *snapshotter) Snapshot(style string, extent *geom.Extent, size image.Point) image.Image {
	a := mbgl.NewLatLng(extent[1], extent[0])
	b := mbgl.NewLatLng(extent[3], extent[2])

	// the mbgl::MapSnapshotter.Snapshot method is not thread safe
	// also, we are making lat/lng changes to the class which cannot be
	// changed through the lifetime of the snapshot routine
	s.snapLock.Lock()

	s.snap.SetStyleURL(style)

	s.size.Width = uint32(size.X)
	s.size.Height = uint32(size.Y)
	s.snap.SetSize(s.size)

	s.snap.SetRegion(mbgl.NewLatLngBounds(a, b))

	_img := s.snap.Snapshot()

	s.snapLock.Unlock()

	img := _img.Image()

	go func() {
		a.Destruct()
		b.Destruct()
		_img.Destruct()
	}()

	return img
}

func SnapshotTile(s Snapshotter, style string, tile slippy.Tile, size image.Point) image.Image {
	ext := tile.Extent4326()
	return s.Snapshot(style, ext, size)
}
