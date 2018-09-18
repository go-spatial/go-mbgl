#ifndef _mbgl_map_h
#define _mbgl_map_h

#include "file_source.h"
#include "map.h"
#include "size.h"
#include "scheduler.h"
#include "renderer_frontend.h"

typedef struct{} MbglMap;

#ifdef __cplusplus
extern "C" {
#endif

MbglMap * mbgl_map_new(MbglRendererFrontend * frontend,
	MbglSize * size,
	float pixelRatio,
	MbglFileSource * src,
	MbglScheduler * sched);

void mbgl_map_destruct(MbglMap * self);

#ifdef __cplusplus
}
#endif

#endif