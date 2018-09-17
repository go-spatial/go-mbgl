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

    std::cout << "in new" << std::endl;

    auto _src = reinterpret_cast<FileSource *>(src);
    auto _sched = reinterpret_cast<Scheduler*>(sched);
    auto _style = std::make_pair((bool) false /*isFile*/, std::string(style));
    auto _size = reinterpret_cast<Size*>(size);
    auto _camOpts = reinterpret_cast<CameraOptions*>(camOpts);
    auto _region = reinterpret_cast<LatLngBounds*>(region);

    std::cout << "variables init" << std::endl;

    optional<std::string> _cacheDir;
    if (cacheDir != nullptr) {
        _cacheDir = std::string(cacheDir);
    }

    MapSnapshotter * ptr;
    try {
        ptr = new MapSnapshotter(
        _src,
        std::shared_ptr<Scheduler>(_sched),
        _style,
        *_size,
        pixelRatio,
        make_optional(_camOpts),
        make_optional(_region),
        _cacheDir);

    } catch (std::runtime_error &e) {
        std::cout << "err" << e.what() << std::endl;
    }

    std::cout << "returning" << std::endl;

    return reinterpret_cast<MbglMapSnapshotter*>(ptr);
}

void mbgl_map_snapshotter_destruct(MbglMapSnapshotter * self) {
    auto cast = reinterpret_cast<MapSnapshotter*>(self);

    std::cout << "del" << std::endl;

    try {
        delete cast;
    } catch (std::runtime_error &e) {
        std::cout << "err" << e.what() << std::endl;
    }
    std::cout << "del-ed" << std::endl;

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
            errStr = std::string("snapshot: ") + std::string(e.what());
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

void mbgl_map_snapshotter_set_style_url(MbglMapSnapshotter * self, const char * style) {
    std::cout << "setting style: " << std::string(style) << std::endl;
    auto ms = reinterpret_cast<MapSnapshotter*>(self);
    ms->setStyleURL(std::string(style));

    std::cout << "set style: " << ms->getStyleURL() << std::endl;
 }

void mbgl_map_snapshotter_set_size(MbglMapSnapshotter * self, MbglSize * size) {
    auto ms = reinterpret_cast<MapSnapshotter*>(self);
    auto _size = reinterpret_cast<Size*>(size);

    ms->setSize(*_size);
}

// image

RawImage * mbgl_premultiplied_image_raw(MbglPremultipliedImage * img) {
    auto _img = reinterpret_cast<PremultipliedImage*>(img);

    auto ret = new RawImage{};
    ret->height = _img->size.height;
    ret->width = _img->size.width;
    ret->data = _img->data.get();

    return ret;
}

void mbgl_premultiplied_image_destruct(MbglPremultipliedImage * self) {
    auto cast = reinterpret_cast<PremultipliedImage*>(self);
    delete cast;
}