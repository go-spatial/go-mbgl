#include <experimental/optional>

#include <mbgl/map/camera.hpp>

#include "camera_options.h"
#include "lat_lng.h"

using namespace mbgl;


template <typename T>
optional<T> toOptional(T * ptr) {
    if (ptr == nullptr) {
        return {};
    }

    return *ptr;
}

// camera options
MbglCameraOptions * mbgl_camera_options_new(MbglLatLng * center,
    MbglEdgeInsets * padding,
    MbglPoint * anchor,
    double * zoom,
    double * angle,
    double * pitch) {

    auto _center = reinterpret_cast<LatLng*>(center);

    EdgeInsets _padding;
    if (padding == nullptr) {
        _padding = EdgeInsets();
    } else {
        _padding = *reinterpret_cast<EdgeInsets*>(padding);
    }

    auto _anchor = reinterpret_cast<ScreenCoordinate*>(anchor);

    auto opts = new CameraOptions{
        toOptional(_center),
        _padding,
        toOptional(_anchor),
        toOptional(zoom),
        toOptional(angle),
        toOptional(pitch)
    };

    return reinterpret_cast<MbglCameraOptions*>(opts);
}

void mbgl_camera_options_update(MbglCameraOptions * self,
    MbglLatLng * center,
    MbglEdgeInsets * padding,
    MbglPoint * anchor,
    double * zoom,
    double * angle,
    double * pitch) {

    auto _center = reinterpret_cast<LatLng*>(center);

    EdgeInsets _padding;
    if (padding == nullptr) {
        _padding = EdgeInsets();
    } else {
        _padding = *reinterpret_cast<EdgeInsets*>(padding);
    }

    auto _anchor = reinterpret_cast<ScreenCoordinate*>(anchor);

    auto opts = reinterpret_cast<CameraOptions*>(self);

    opts->center = toOptional(_center);
    opts->padding = _padding;
    opts->anchor = toOptional(_anchor);
    opts->zoom = toOptional(zoom);
    opts->angle = toOptional(angle);
    opts->pitch = toOptional(pitch);
}

void mbgl_camera_options_destruct(MbglCameraOptions * self) {
    auto cast = reinterpret_cast<CameraOptions*>(self);
    delete self;
}

// edge insets
MbglEdgeInsets * mbgl_edge_insets_new(double top, double left, double bottom, double right) {
    auto edges = new EdgeInsets(top, left, bottom, right);

    return reinterpret_cast<MbglEdgeInsets*>(edges);
}

void mbgl_edge_insets_destruct(MbglEdgeInsets * self) {
    auto cast = reinterpret_cast<EdgeInsets*>(self);

    delete cast;
}

// point
MbglPoint * mbgl_point_new(double x, double y) {
    auto pt = new ScreenCoordinate{
        x, y,
    };

    return reinterpret_cast<MbglPoint*>(pt);
}

void mbgl_point_update(MbglPoint * self, double x, double y) {
    auto pt = reinterpret_cast<ScreenCoordinate*> (self);

    pt->x = x;
    pt->y = y;
}

void mbgl_point_delete(MbglPoint * self) {
    auto cast = reinterpret_cast<ScreenCoordinate*>(self);
    delete cast;
}