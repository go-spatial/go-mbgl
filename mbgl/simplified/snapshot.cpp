
#include <iostream>
#include "snapshot.h"
#include <mbgl/map/map.hpp>
#include <mbgl/util/image.hpp>
#include <mbgl/util/run_loop.hpp>
#include <mbgl/util/default_styles.hpp>

#include <mbgl/gl/headless_frontend.hpp>
#include <mbgl/util/default_thread_pool.hpp>
#include <mbgl/storage/default_file_source.hpp>
#include <mbgl/style/style.hpp>

snapshot_Result Snapshot(snapshot_Params params) {

    snapshot_Result result; 

    result.DidError = 0;
    result.Image    = NULL;
    result.Err      = (char *)("");

    mbgl::ThreadPool threadPool(4);

    mbgl::DefaultFileSource fileSource(params.cache_file, params.asset_root);
    mbgl::HeadlessFrontend frontend({ params.width, params.height }, float(params.ppi_ratio), fileSource, threadPool);

    mbgl::Map map(frontend, mbgl::MapObserver::nullObserver(), frontend.getSize(), params.ppi_ratio, fileSource, threadPool, mbgl::MapMode::Static);

    map.getStyle().loadURL(params.style);
    map.setLatLngZoom({ params.lat, params.lng }, params.zoom);
    map.setBearing(params.bearing);
    map.setPitch(params.pitch);

    try {

       auto img1 = frontend.render(map);
	    auto _img = new mbgl::PremultipliedImage();
	    *_img = std::move(img1);
		 auto data = _img->data.get();
       result.Image = new RawImage{};
       result.Image->Height = _img->size.height;
       result.Image->Width  = _img->size.width;
       result.Image->Data   = data;
       result.DidError     = 0;
       result.Err          = (char *)("");

    } catch(std::exception& e) {

       result.Image    = NULL;
       result.DidError = 1;
       result.Err      = e.what();

    }
    return result;
}

