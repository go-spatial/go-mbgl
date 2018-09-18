#ifndef _mbgl_headless_frontend_h
#define _mbgl_headless_frontend_h


#include "file_source.h"
#include "map_snapshotter.h"
#include "size.h"
#include "scheduler.h"
#include "map.h"


typedef struct {} MbglHeadlessFrontend;

#ifdef __cplusplus
extern "C" {
#endif

MbglHeadlessFrontend * mbgl_headless_frontend_new(
	MbglSize * size,
	float pixelRatio,
	MbglFileSource * source,
	MbglScheduler * sched,
	const char * cacheDir,
	const char * fontFamily);

void mbgl_headless_frontend_destruct(MbglHeadlessFrontend * self);

MbglPremultipliedImage * mbgl_headless_frontend_render(MbglHeadlessFrontend * self, MbglMap * map);

MbglSize * mbgl_headless_frontend_get_size();
void mbgl_headless_frontend_set_size(MbglSize *);


#ifdef __cplusplus
}
#endif

#endif