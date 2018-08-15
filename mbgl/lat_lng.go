package mbgl

/*
#include "lat_lng.h"
 */
import "C"
import "github.com/go-spatial/geom/slippy"

//TODO (@ear7h): use github.com/go-spatial/geom type
type LatLng struct {
	ptr *C.MbglLatLng
}

func NewLatLng(lat, lng float64) *LatLng {
	ptr := C.mbgl_lat_lng_new(C.double(lat), C.double(lng))

	return &LatLng{ptr:ptr}
}

func (ll *LatLng) Destruct() {
	C.mbgl_lat_lng_destruct(ll.ptr)
	ll.ptr = nil
}


// bounds

type LatLngBounds struct {
	ptr *C.MbglLatLngBounds
}

func NewLatLngBounds(a, b *LatLng) *LatLngBounds {
	ptr := C.mbgl_lat_lng_bounds_hull(a.ptr, b.ptr)
	return &LatLngBounds{ptr:ptr}
}

func NewLatLngBoundsFromTile(tile *slippy.Tile) *LatLngBounds {
	// east, north, west, south
	bound := tile.Bounds()
	a := NewLatLng(bound[1], bound[0])
	b := NewLatLng(bound[3], bound[2])

	ret := NewLatLngBounds(a, b)

	a.Destruct()
	b.Destruct()

	return ret
}

func (bb *LatLngBounds) Destruct() {
	C.mbgl_lat_lng_bounds_destruct(bb.ptr)
}