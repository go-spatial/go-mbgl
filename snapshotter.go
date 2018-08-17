package mbgl

import (
	"image"
	"runtime"
	"sync"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/go-mbgl/mbgl"
	"github.com/go-spatial/geom/slippy"
)

// This will set the size of the mbgl::ThreadPool
var DefaultThreadPoolSize = 4

// Snapshotter is an interface for processing raster image from
// vector tiles.
type Snapshotter interface {
	SetStyle(style string)
	Snapshot(extent *geom.Extent, size image.Point) image.Image
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


// High level function for taking a single snappshot. This simply creates a new snapshotter and uses it.
func Snapshot(style string, ext *geom.Extent, size image.Point) image.Image {
	return NewSnapshotter(style).Snapshot(ext, size)
}

// This creates an instance of a Snapshotter with the specified style.
// Note: this high level implementation is thread safe, but performance might be better to lower the DefaultThreadPoolSize and use multiple snapshotters
// TODO(@ear7h): write benchmarks
func NewSnapshotter(style string) Snapshotter {
	src := mbgl.NewDefaultFileSource("", "", nil)
	runtime.SetFinalizer(src, (*mbgl.DefaultFileSource).Destruct)

	tpool := mbgl.NewThreadPool(DefaultThreadPoolSize)
	mbgl.SchedulerSetCurrent(tpool)

	size := mbgl.Size{Width: 100, Height: 100}

	snap := mbgl.NewMapSnapshotter(src,
		tpool,
		style,
		size,
		1.0,
		nil,
		nil,
		nil)

	ret := &snapshotter{
		fsrc: src,
		snap: snap,
		size: size,
	}

	runtime.SetFinalizer(ret, (*snapshotter).destruct)

	return ret
}

// Take the snapshot with a *geom.Extent encoded as lat/lng (ie. WSG84).
func (s *snapshotter) Snapshot(extent *geom.Extent, size image.Point) image.Image {
	a := mbgl.NewLatLng(extent[1], extent[0])
	b := mbgl.NewLatLng(extent[3], extent[2])

	// the mbgl::MapSnapshotter.Snapshot method is not thread safe
	// also, we are making lat/lng changes to the class which cannot be
	// changed through the lifetime of the snapshot routine
	s.snapLock.Lock()

	s.snap.SetRegion(mbgl.NewLatLngBounds(a, b))

	s.size.Width = uint32(size.X)
	s.size.Height = uint32(size.Y)
	s.snap.SetSize(s.size)

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

func (s *snapshotter) SetStyle(style string) {
	s.snap.SetStyleURL(style)
}

func SnapshotTile(s Snapshotter, tile slippy.Tile, size image.Point) image.Image {
	ext := tile.Extent(geom.WGS84)
	return s.Snapshot(ext, size)
}

