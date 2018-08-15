#include <iostream>

#include <mbgl/map/map_snapshotter.hpp>
#include <mbgl/map/camera.hpp>
#include <mbgl/storage/file_source.hpp>
#include <mbgl/actor/scheduler.hpp>

#include "map_snapshotter.h"
#include "file_source.h"
#include "scheduler.h"
#include "size.h"
#include "camera_options.h"
#include "lat_lng.h"

using namespace mbgl;

template <typename T>
optional<T> make_optional(T *t) {
    if (t == nullptr) {
        return {};
    }

    return *t;
}

MbglMapSnapshotter * mbgl_map_snapshotter_new(
    MbglFileSource * src,
    MbglScheduler * sched,
    int isFile, const char * style,
    MbglSize * size,
    float pixelRatio,
    MbglCameraOptions * camOpts,
    MbglLatLngBounds * region,
    const char * cacheDir) {

    auto _src = reinterpret_cast<FileSource *>(src);
    auto _sched = reinterpret_cast<Scheduler*>(sched);
    auto _style = std::make_pair((bool) isFile, std::string(style));
    auto _size = reinterpret_cast<Size*>(size);
    auto _camOpts = reinterpret_cast<CameraOptions*>(camOpts);
    auto _region = reinterpret_cast<LatLngBounds*>(region);

    optional<std::string> _cacheDir;
    if (cacheDir != nullptr) {
        _cacheDir = std::string(cacheDir);
    }

    auto ptr = new MapSnapshotter(
        _src,
        std::shared_ptr<Scheduler>(_sched),
        _style,
        *_size,
        pixelRatio,
        make_optional(_camOpts),
        make_optional(_region),
        _cacheDir);

    return reinterpret_cast<MbglMapSnapshotter*>(ptr);
}

void mbgl_map_snapshotter_destruct(MbglMapSnapshotter * self) {
    auto cast = reinterpret_cast<MapSnapshotter*>(self);

    delete cast;
}

MbglPremultipliedImage * mbgl_map_snapshotter_snapshot(MbglMapSnapshotter * self) {
    auto ms = reinterpret_cast<MapSnapshotter*>(self);

    PremultipliedImage img1;
    std::exception_ptr err1;

    auto cb = std::make_unique<mbgl::Actor<mbgl::MapSnapshotter::Callback>>(
        *mbgl::Scheduler::GetCurrent(),
        [&img1, &err1] (std::exception_ptr err,
            mbgl::PremultipliedImage img,
            mbgl::MapSnapshotter::Attributions attr,
            mbgl::MapSnapshotter::PointForFn pt,
            mbgl::MapSnapshotter::LatLngForFn ll) {
                // std::cout << "CALLBACK" << std::endl;
                err1 = std::move(err);
                img1 = std::move(img);

//                if (err1) {
//                    try {
//                        std::rethrow_exception(err1);
//                    } catch (std::exception &e) {
//                        std::cout << "err" << e.what() << std::endl;
//                    }
//                }

                });

    ms->snapshot(cb->self());


    // run the event loop until the image is finished processing
    while(!img1.valid() && !err1) {
        mbgl::util::RunLoop::Get()->runOnce();
    }


    std::string errStr;
    if (err1) {
        try {
            std::rethrow_exception(err1);
        } catch (std::exception &e) {
            std::cout << "err" << e.what() << std::endl;
            errStr = e.what();
        }
    }

    auto heap = new PremultipliedImage();
    *heap = std::move(img1);

    return reinterpret_cast<MbglPremultipliedImage*>(heap);
}

void mbgl_map_snapshotter_set_camera_options(MbglMapSnapshotter * self, MbglCameraOptions * camOpts) {
    auto ms = reinterpret_cast<MapSnapshotter*>(self);
    auto _camOpts = reinterpret_cast<CameraOptions*>(camOpts);

    ms->setCameraOptions(*_camOpts);
}

void mbgl_map_snapshotter_set_region(MbglMapSnapshotter * self, MbglLatLngBounds * region) {
    auto ms = reinterpret_cast<MapSnapshotter*>(self);
    auto _region = reinterpret_cast<LatLngBounds*>(region);

    ms->setRegion(*_region);
}

RawImage * mbgl_premultiplied_image_raw(MbglPremultipliedImage * img) {
    auto _img = reinterpret_cast<PremultipliedImage*>(img);

    auto ret = new RawImage{};
    ret->height = _img->size.height;
    ret->width = _img->size.width;
    ret->data = _img->data.get();

    return ret;
}

void mbgl_premultiplied_image_raw_delete(MbglPremultipliedImage * self) {
    auto cast = reinterpret_cast<PremultipliedImage*>(self);

    delete cast;
}