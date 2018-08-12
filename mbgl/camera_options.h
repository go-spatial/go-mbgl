#ifndef _mbgl_camera_options_h
#define _mbgl_camera_options_h

#include "lat_lng.h"

typedef struct{} MbglCameraOptions;
typedef struct{} MbglEdgeInsets;
typedef struct{} MbglPoint;


#ifdef __cplusplus
extern "C"{
#endif

// camera options
MbglCameraOptions * mbgl_camera_options_new(MbglLatLng * latLng,
    MbglEdgeInsets * padding,
    MbglPoint * anchor,
    double * zoom,
    double * angle,
    double * pitch);

void mbgl_camera_options_update(MbglCameraOptions * self,
    MbglLatLng * latLng,
    MbglEdgeInsets * padding,
    MbglPoint * anchor,
    double * zoom,
    double * angle,
    double * pitch);

void mbgl_camera_options_destruct(MbglCameraOptions * self);

// edge insets
MbglEdgeInsets * mbgl_edge_insets_new(double top, double left, double bottom, double right);
void mbgl_edge_insets_destruct(MbglEdgeInsets * self);

// point
MbglPoint * mbgl_point_new(double x, double y);
void mbgl_point_update(MbglPoint * self, double x, double y);
void mbgl_point_destruct(MbglPoint * self);

#ifdef __cplusplus
} //extern "C"
#endif


#endif // _mbgl_camera_options_h