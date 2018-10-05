#ifndef _mbgl_simplifed_snapshot_h
#define _mbgl_simplifed_snapshot_h

#include <stdlib.h>
#include <stdint.h>


// This is the raw image that we be returned.
typedef struct{
	size_t Height;
	size_t Width;
	uint8_t *Data;
} RawImage;

typedef struct {
	char * style;
	char * cache_file;
	char * asset_root;
	uint32_t width;
	uint32_t height;
	double ppi_ratio;
	double lat;
	double lng;
	double zoom;
	double pitch;
	double bearing;
} snapshot_Params;

typedef struct {
       RawImage * Image;
       int DidError;
       const char * Err;
} snapshot_Result;


#ifdef __cplusplus
extern "C" {
#endif

  snapshot_Result Snapshot(snapshot_Params params);

#ifdef __cplusplus
} // extern "C"
#endif

#endif

