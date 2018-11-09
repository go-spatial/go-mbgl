/*
bounds

This is a temporary package to implement BoundsCenterZoom
We are hardcodeing thing here till the geom and proj packages are ready
can support these functions.
*/

package bounds

import (
	"math"

	"github.com/go-spatial/geom"
)

type matTransform struct {
	xscale     float64
	xtranslate float64
	yscale     float64
	ytranslate float64
}

func (mt *matTransform) Transform(pt [2]float64, scale float64) [2]float64 {

	x, y := pt[0], pt[1]
	if scale == 0.0 {
		scale = 1.0
	}
	return [2]float64{
		scale * ((mt.xscale * x) + mt.xtranslate),
		scale * ((mt.yscale * y) + mt.ytranslate),
	}
}

func (mt *matTransform) Untransform(pt [2]float64, scale float64) [2]float64 {
	x, y := pt[0], pt[1]
	if scale == 0.0 {
		scale = 1.0
	}
	return [2]float64{
		(x/scale - mt.xtranslate) / mt.xscale,
		(y/scale - mt.ytranslate) / mt.yscale,
	}
}

var projections = [...]struct {
	name               string
	radius             float64
	maxLatitude        float64
	circumferenceRatio float64
	tTranslate         float64
	bounds             *geom.Extent
	transformer        matTransform
}{
	{
		name:               "ESPG3857",
		radius:             6378137, // earth radius for ESPG3857
		circumferenceRatio: 1 / (2 * math.Pi * 6378137),
		maxLatitude:        85.0511287798,
		tTranslate:         0.5,
		bounds:             &geom.Extent{-180.0, -85.06, 180.0, 85.06},
	},
}

func init() {
	for i := range projections {
		prj := projections[i]
		projections[i].transformer = matTransform{
			xscale:     prj.circumferenceRatio,
			xtranslate: prj.tTranslate,
			yscale:     -prj.circumferenceRatio,
			ytranslate: prj.tTranslate,
		}
	}
}

type AProjection int

const (
	ESPG3857 = AProjection(0)
)

func (p AProjection) String() string       { return projections[int(p)].name }
func (p AProjection) Bounds() *geom.Extent { return projections[int(p)].bounds }
func (p AProjection) R() float64           { return projections[int(p)].radius }
func (p AProjection) MaxLatitude() float64 { return projections[int(p)].maxLatitude }

func (p AProjection) Transform(pt [2]float64, scale float64) [2]float64 {
	return projections[int(p)].transformer.Transform(pt, scale)
}
func (p AProjection) Untransform(pt [2]float64, scale float64) [2]float64 {
	return projections[int(p)].transformer.Untransform(pt, scale)
}

func (p AProjection) Project(latlng [2]float64) (xy [2]float64) {
	lat, lng := latlng[0], latlng[1]
	d := math.Pi / 180
	max := p.MaxLatitude()
	r := p.R()
	_lat := math.Max(math.Min(max, lat), -max)
	sin := math.Sin(_lat * d)

	return [2]float64{r * lng * d, r * math.Log((1+sin)/(1-sin)) / 2}
}

func (p AProjection) Unproject(pt [2]float64) (latlng [2]float64) {
	d := 180 / math.Pi
	prj := projections[p]

	return [2]float64{
		(2*math.Atan(math.Exp(pt[1]/prj.radius)) - (math.Pi / 2)) * d,
		pt[0] * d / prj.radius,
	}
}

func ZoomTile(bounds *geom.Extent, width, height float64, tileSize int) float64 {
	// assume ESPG3857 for now.
	prj := ESPG3857
	if bounds == nil {
		// we want the whole world.
		bounds = prj.Bounds()
	}

	// for lat lng geom.Extent should be laid out as follows:
	// {west, south, east, north}
	nw := [2]float64{bounds[3], bounds[0]}
	se := [2]float64{bounds[1], bounds[2]}

	// 256 is the tile size.
	ptupper := LatLngToPoint(prj, nw[0], nw[1], 0, tileSize)
	ptlower := LatLngToPoint(prj, se[0], se[1], 0, tileSize)

	b := geom.NewExtent(ptupper, ptlower)
	scale := math.Min(width/b.XSpan(), height/b.YSpan())
	return math.Floor(math.Log(scale) / math.Ln2)
}

//	Zoom returns the zoom level for supplied bounds
//	useful when rendering static map images
// tile size is assumed to be 256
//
//	TODO: add padding support
func Zoom(bounds *geom.Extent, width, height float64) float64 {
	return ZoomTile(bounds, width, height, 256)
}

func CenterTile(bounds *geom.Extent, zoom float64, tileSize int) [2]float64 {
	// assume ESPG3857 for now.
	prj := ESPG3857
	if bounds == nil {
		// we want the whole world.
		bounds = prj.Bounds()
	}

	// for lat lng geom.Extent should be laid out as follows:
	// {west, south, east, north}
	ne := [2]float64{bounds[3], bounds[2]}
	sw := [2]float64{bounds[1], bounds[0]}

	swPt := LatLngToPoint(prj, sw[0], sw[1], zoom, tileSize)
	nePt := LatLngToPoint(prj, ne[0], ne[1], zoom, tileSize)

	// center point.
	centerPtX := (swPt[0] + nePt[0]) / 2
	centerPtY := (swPt[1] + nePt[1]) / 2

	// 256 is the tile size.
	lat, lng := PointToLatLng(prj, [2]float64{centerPtX, centerPtY}, zoom, tileSize)
	return [2]float64{lat, lng}
}

func WidthHeightTile(bounds *geom.Extent, zoom float64, tileSize int) (width, height float64) {
	// assume ESPG3857 for now.
	prj := ESPG3857
	if bounds == nil {
		// we want the whole world.
		bounds = prj.Bounds()
	}

	// for lat lng geom.Extent should be laid out as follows:
	// {west, south, east, north}
	ne := [2]float64{bounds[3], bounds[2]}
	sw := [2]float64{bounds[1], bounds[0]}

	swPt := LatLngToPoint(prj, sw[0], sw[1], zoom, tileSize)
	nePt := LatLngToPoint(prj, ne[0], ne[1], zoom, tileSize)

	return math.Abs(swPt[0] - nePt[0]), math.Abs(swPt[1] - nePt[1])
}
func WidthHeight(bounds *geom.Extent, zoom float64) (float64, float64) {
	return WidthHeightTile(bounds, zoom, 256)
}

func Center(bounds *geom.Extent, zoom float64) [2]float64 {
	return CenterTile(bounds, zoom, 256)
}

func CenterZoomTile(bounds *geom.Extent, width, height float64, tileSize int) ([2]float64, float64) {
	zoom := ZoomTile(bounds, width, height, tileSize)
	return CenterTile(bounds, zoom, tileSize), zoom
}

func CenterZoom(bounds *geom.Extent, width, height float64) ([2]float64, float64) {
	return CenterZoomTile(bounds, width, height, 256)
}

func ScaleTile(zoom float64, tileSize int) float64 {
	return float64(tileSize) * math.Pow(2, zoom)
}

func Scale(zoom float64) float64 { return ScaleTile(zoom, 256) }

// type AProjection int

func LatLngToPoint(prj AProjection, lat, lng, zoom float64, tilesize int) [2]float64 {
	prjPt := prj.Project([2]float64{lat, lng})
	scale := ScaleTile(zoom, tilesize)
	return prj.Transform(prjPt, scale)
}

func PointToLatLng(prj AProjection, point [2]float64, zoom float64, tilesize int) (lat, lng float64) {
	scale := ScaleTile(zoom, tilesize)
	utPt := prj.Untransform(point, scale)
	latlng := prj.Unproject(utPt)
	return latlng[0], latlng[1]
}
