package mbgl
import "C"

/*
#include "camera_options.h"
 */
import "C"

type EdgeInsets C.MbglEdgeInsets

func (ei *EdgeInsets) edgeInsets() *C.MbglEdgeInsets {
	return (*C.MbglEdgeInsets)(ei)
}

func NewEdgeInsets(top, left, bottom, right float64) *EdgeInsets {
	ptr := C.mbgl_edge_insets_new(
		C.double(top),
		C.double(left),
		C.double(bottom),
		C.double(right),
	)

	return (*EdgeInsets)(ptr)
}

func (ei *EdgeInsets) Destruct() {
	C.mbgl_edge_insets_destruct(ei.edgeInsets())
}

type Point struct {
	X, Y float64

	ptr *C.MbglPoint
}

func (p *Point) update() {
	if p.ptr == nil {
		ptr := C.mbgl_point_new(
			C.double(p.X),
			C.double(p.Y))

		p.ptr = ptr
	} else {
		C.mbgl_point_update(
			p.ptr,
			C.double(p.X),
			C.double(p.Y))
	}
}

func (p *Point) cPtr() *C.MbglPoint {
	if p == nil {return nil}

	p.update()
	return p.ptr
}

type CameraOptions struct {
	Center *LatLng
	Padding *EdgeInsets
	Anchor *Point
	Zoom *float64
	Angle *float64
	Pitch *float64

	ptr *C.MbglCameraOptions
}

func (opt *CameraOptions) update() {

	// todo (@ear7h): change structs to wrapped types
	var center *C.MbglLatLng
	if opt.Center != nil {
		center = opt.Center.latLng()
	}

	var padding *C.MbglEdgeInsets
	if opt.Padding != nil {
		padding = opt.Padding.edgeInsets()
	}

	if opt.ptr == nil {
		ptr := C.mbgl_camera_options_new(
			center,
			padding,
			opt.Anchor.cPtr(),
			(*C.double)(opt.Zoom),
			(*C.double)(opt.Angle),
			(*C.double)(opt.Pitch))
		opt.ptr = ptr
	} else {
		C.mbgl_camera_options_update(
			opt.ptr,
			center,
			padding,
			opt.Anchor.cPtr(),
			(*C.double)(opt.Zoom),
			(*C.double)(opt.Angle),
			(*C.double)(opt.Pitch))
	}
}

func (opt *CameraOptions) cPtr() *C.MbglCameraOptions {
	if opt == nil {return nil}

	opt.update()
	return opt.ptr
}

func (opt *CameraOptions) Destruct() {
	C.mbgl_camera_options_destruct(opt.ptr)
	opt.ptr = nil
}
