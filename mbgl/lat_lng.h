#ifndef _mbgl_lat_lng_h
#define _mbgl_lat_lng_h

typedef struct{} MbglLatLng;
typedef struct{} MbglLatLngBounds;

#ifdef __cplusplus
extern "C"{
#endif

// lat long
MbglLatLng * mbgl_lat_lng_new(double lat, double lng);
void mbgl_lat_lng_destruct(MbglLatLng * self);

// bounds
MbglLatLngBounds * mbgl_lat_lng_bounds_hull(MbglLatLng * a, MbglLatLng * b);
void mbgl_lat_lng_bounds_destruct(MbglLatLngBounds * self);

#ifdef __cplusplus
} //extern "C"
#endif

#endif //_mbgl_lat_lng_h