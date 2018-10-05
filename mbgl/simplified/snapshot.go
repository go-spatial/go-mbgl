package simplified

import (
	"errors"
	"log"
	"runtime"
	"strings"
	"unsafe"
)

/*
#include <stdio.h>
#include "snapshot.h"


snapshot_Params NewSnapshotParams( char* style, char* cacheFile, char* assetRoot, uint32_t width, uint32_t height, int ppi_ratio, double lat, double lng, double zoom, double pitch, double bearing) {

	snapshot_Params params;
	params.style      = style;
	params.cache_file = cacheFile;
	params.asset_root = assetRoot;
	params.width      = width;
	params.height     = height;
	params.ppi_ratio  = ppi_ratio;
	params.lat        = lat;
	params.lng        = lng;
	params.zoom       = zoom;
	params.pitch      = pitch;
	params.bearing    = bearing;
	return params;
}
*/
import "C"

type Snapshotter struct {
	Style     string
	CacheFile string
	AssetRoot string
	Width     uint32
	Height    uint32
	PPIRatio  float64
	Lat       float64
	Lng       float64
	Zoom      float64
	Pitch     float64
	Bearing   float64
}

func (snap Snapshotter) AsParams() (p C.snapshot_Params, err error) {
	style := snap.Style
	if strings.TrimSpace(style) == "" {
		return p, errors.New("Style param is required.")
	}
	cacheFile := snap.CacheFile
	if strings.TrimSpace(cacheFile) == "" {
		cacheFile = "cache.sql"
	}
	assetRoot := snap.AssetRoot
	if strings.TrimSpace(assetRoot) == "" {
		assetRoot = "."
	}
	width := snap.Width
	if width == 0 {
		width = 512
	}
	height := snap.Height
	if height == 0 {
		height = 512
	}
	zoom := snap.Zoom
	if zoom <= 0 {
		zoom = 0
	}
	ppiratio = snap.PPIRatio
	if ppiratio == 0.0 {
		ppiratio = 1.0
	}
	img := C.NewSnapshotParams(
		C.CString(style),
		C.CString(cacheFile),
		C.CString(assetRoot),
		C.uint32_t(width),
		C.uint32_t(height),
		C.double(ppiratio),
		C.double(snap.Lat),
		C.double(snap.Lng),
		C.double(zoom),
		C.double(snap.Pitch),
		C.double(snap.Bearing),
	)
	return img, nil
}

func Snapshot(snap Snapshotter) (img Image, err error) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	_params, err := snap.AsParams()
	if err != nil {
		return img, err
	}

	result := C.Snapshot(_params)
	didError := (int)(result.DidError)
	if didError == 1 {
		return img, errors.New(C.GoString(result.Err))
	}

	img.Width, img.Height = int(result.Image.Width), int(result.Image.Height)
	bytes := img.Width * img.Height * 4
	log.Printf("Width %v : Height %v : Bytes %v", img.Width, img.Height, bytes)
	img.Data = C.GoBytes(unsafe.Pointer(result.Image.Data), C.int(bytes))

	return img, nil
}
