package simplified

import (
	"context"
	"errors"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/*
#include <stdio.h>
#include "snapshot.h"


snapshot_Params NewSnapshotParams( char* style, char* cacheFile, char* assetRoot, uint32_t width, uint32_t height, double ppi_ratio, double lat, double lng, double zoom, double pitch, double bearing) {

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

var ErrManagerExiting = errors.New("Manager shutting down.")

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
	ppiratio := snap.PPIRatio
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

func snapshot(snap Snapshotter) (img Image, err error) {

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
	data := C.GoBytes(unsafe.Pointer(result.Image.Data), C.int(bytes))
	img.Data = make([]byte, len(data))
	copy(img.Data, data)
	C.free(unsafe.Pointer(result.Image.Data))

	return img, nil
}

type snapJobReply struct {
	image Image
	err   error
}

type snapJob struct {
	Params Snapshotter
	Reply  chan<- snapJobReply
}

func snapshotWorker(job snapJob) {
	img, err := snapshot(job.Params)
	job.Reply <- snapJobReply{
		image: img,
		err:   err,
	}
}

var rwManagerRunning sync.RWMutex
var isManagerRunning bool
var workQueue chan snapJob

func IsManagerRunning() (b bool) {
	rwManagerRunning.RLock()
	defer rwManagerRunning.RUnlock()
	return isManagerRunning
}

// StartSnapshotManager will block till the Manager has started.
func StartSnapshotManager(ctx context.Context) {
	go SnapshotManager(ctx)
	for {
		<-time.After(10 * time.Millisecond)
		if IsManagerRunning() {
			return
		}
	}
}

func SnapshotManager(ctx context.Context) {
	if IsManagerRunning() {
		return
	}

	log.Println("Starting up manager.")
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	NewRunLoop()
	defer DestroyRunLoop()
	workQueue = make(chan snapJob)

	rwManagerRunning.Lock()
	isManagerRunning = true
	rwManagerRunning.Unlock()
	defer func() {
		rwManagerRunning.Lock()
		isManagerRunning = false
		rwManagerRunning.Unlock()
	}()

	for {
		select {
		case job := <-workQueue:
			snapshotWorker(job)
		case <-ctx.Done():
			// We are exiting...
			rwManagerRunning.Lock()
			isManagerRunning = false
			rwManagerRunning.Unlock()
			break
		}
	}
	for j := range workQueue {
		// We are shutdowning, let's empty the queue.
		j.Reply <- snapJobReply{
			err: ErrManagerExiting,
		}
	}
}

func Snapshot1(snap Snapshotter) (img Image, err error) {
	if !IsManagerRunning() {
		if workQueue != nil {
			close(workQueue)
		}
		return img, ErrManagerExiting
	}
	var reply = make(chan snapJobReply)

	workQueue <- snapJob{
		Params: snap,
		Reply:  reply,
	}

	r := <-reply
	close(reply)
	return r.image, r.err
}
