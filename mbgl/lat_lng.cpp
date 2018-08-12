#include <mbgl/util/geo.hpp>

#include "lat_lng.h"

using namespace mbgl;

//lat long
MbglLatLng * mbgl_lat_lng_new(double lat, double lng) {
    auto latLng = new LatLng{
        lat, lng
    };

    return reinterpret_cast<MbglLatLng*>(latLng);
}

void mbgl_lat_lng_destruct(MbglLatLng * self) {
    auto latLng = reinterpret_cast<LatLng*>(self);

    delete latLng;
}

// bounds
MbglLatLngBounds * mbgl_lat_lng_bounds_hull(MbglLatLng * a, MbglLatLng * b) {
    auto _a = reinterpret_cast<LatLng*>(a);
    auto _b = reinterpret_cast<LatLng*>(b);

    auto bb = LatLngBounds::hull(*_a, *_b);
    auto _bb =  new LatLngBounds(bb);

    return reinterpret_cast<MbglLatLngBounds*>(_bb);
}

void mbgl_lat_lng_bounds_destruct(MbglLatLngBounds * self) {
    auto cast = reinterpret_cast<LatLngBounds*>(self);

    delete cast;
}