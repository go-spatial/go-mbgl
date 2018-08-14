#ifndef _mbgl_map_snapshotter_h
#define _mbgl_map_snapshotter_h

#include <stdlib.h>

#include "file_source.h"
#include "scheduler.h"
#include "camera_options.h"
#include "lat_lng.h"
#include "size.h"

typedef struct{} MbglMapSnapshotter;
typedef struct{} MbglPremultipliedImage;

typedef struct{
    size_t height;
    size_t width;
    uint8_t * data;
} RawImage;

#ifdef __cplusplus
extern "C" {
#endif

MbglMapSnapshotter * mbgl_map_snapshotter_new(
    MbglFileSource * src,
    MbglScheduler * sched,
    int isFile, const char * style,
    MbglSize * size,
    float pixelRatio,
    MbglCameraOptions * camOpts,
    MbglLatLngBounds * region,
    const char * cacheDir);

MbglPremultipliedImage * mbgl_map_snapshotter_snapshot(MbglMapSnapshotter * self);

void mbgl_map_snapshotter_set_camera_options(MbglMapSnapshotter * self, MbglCameraOptions * camOpts);

void mbgl_map_snapshotter_destruct(MbglMapSnapshotter * self);

// image
RawImage * mbgl_premultiplied_image_raw(MbglPremultipliedImage * self);

void mbgl_premultiplied_image_destruct(MbglPremultipliedImage * self);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // _mbgl_map_snapshotter_h