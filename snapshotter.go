package mbgl

import (
	"image"
	"runtime"
	"sync"

	"github.com/go-spatial/geom"
	"github.com/go-spatial/go-mbgl/mbgl"
	"github.com/go-spatial/geom/slippy"
)

type Snapshotter interface {
	Snapshot(extent *geom.Extent, size image.Point) image.Image
}

type snapshotter struct {
	fsrc mbgl.FileSource
	snap *mbgl.MapSnapshotter
	size mbgl.Size

	snapLock sync.Mutex
}

var defaultSnapshotter = NewSnapshotter("").(*snapshotter)
var defaultSnapshotterLock = sync.Mutex{}

func Snapshot(src string, ext *geom.Extent, size image.Point) image.Image {
	defaultSnapshotterLock.Lock()
	defer defaultSnapshotterLock.Unlock()

	defaultSnapshotter.snap.SetStyleURL(src)
	return defaultSnapshotter.Snapshot(ext, size)
}

func NewSnapshotter(style string) Snapshotter {
	src := mbgl.NewDefaultFileSource("", "", nil)
	runtime.SetFinalizer(src, (*mbgl.DefaultFileSource).Destruct)

	tpool := mbgl.NewThreadPool(4)
	mbgl.SchedulerSetCurrent(tpool)

	size := mbgl.Size{Width: 100, Height: 100}

	ret := mbgl.NewMapSnapshotter(src,
		tpool,
		style,
		size,
		1.0,
		nil,
		nil,
		nil)

	runtime.SetFinalizer(ret, (*mbgl.MapSnapshotter).Destruct)

	return &snapshotter{
		fsrc: src,
		snap: ret,
		size: size,
	}
}

func (s *snapshotter) Snapshot(extent *geom.Extent, size image.Point) image.Image {
	s.snapLock.Lock()
	defer s.snapLock.Unlock()

	a := mbgl.NewLatLng(extent[1], extent[0])
	b := mbgl.NewLatLng(extent[3], extent[2])

	s.snap.SetRegion(mbgl.NewLatLngBounds(a, b))

	a.Destruct()
	b.Destruct()

	s.size.Width = uint32(size.X)
	s.size.Height = uint32(size.Y)
	s.snap.SetSize(s.size)

	_img := s.snap.Snapshot()
	img := _img.Image()
	_img.Destruct()

	return img
}

func (s *snapshotter) SnapshotTile(tile slippy.Tile, size image.Point) image.Image {
	ext := tile.Extent(geom.WGS84)
	return s.Snapshot(ext, size)
}
