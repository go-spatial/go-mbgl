#include <mbgl/map/map.hpp>
#include <mbgl/util/image.hpp>
#include <mbgl/util/run_loop.hpp>
#include <mbgl/util/default_styles.hpp>

#include <mbgl/gl/headless_frontend.hpp>
#include <mbgl/util/default_thread_pool.hpp>
#include <mbgl/storage/default_file_source.hpp>
#include <mbgl/style/style.hpp>


#include <cstdlib>
#include <iostream>
#include <fstream>


struct snapshot_params {
	std::string style;
	std::string cache_file;
	std::string asset_root;
	uint32_t width;
	uint32_t height;
	int ppi_ratio;
	double lat;
	double lng;
	double zoom;
	double pitch;
	double bearing;
};

struct snapshot_result {
	mbgl::PremultipliedImage image;
	bool did_error;
	std::string err;
};


struct snapshot_result snapshot( struct snapshot_params params ) {

    struct snapshot_result result; 
    result.did_error = false;

    mbgl::ThreadPool threadPool(4);

    mbgl::DefaultFileSource fileSource(params.cache_file, params.asset_root);
    mbgl::HeadlessFrontend frontend({ params.height, params.width }, params.ppi_ratio, fileSource, threadPool);

    mbgl::Map map(frontend, mbgl::MapObserver::nullObserver(), frontend.getSize(), params.ppi_ratio, fileSource, threadPool, mbgl::MapMode::Static);

    //map.getStyle().loadURL("file://style.json");
    map.getStyle().loadURL(params.style);
    map.setLatLngZoom({ params.lat, params.lng }, params.zoom);
    map.setBearing(params.bearing);
    map.setPitch(params.pitch);

    try {
	result.image = frontend.render(map);
    } catch(std::exception& e) {
	result.did_error = true;
	result.err = e.what();
    }
    return result;
}


int main(int argc __attribute__((unused)), char *argv[] __attribute__((unused)) ) {
    const std::string output     = "out.png";

    using namespace mbgl;

    util::RunLoop loop;

    struct snapshot_params snpParams;
    snpParams.cache_file = "cache.sqlite";
    snpParams.asset_root = ".";
    snpParams.style = "file://style.json";
    snpParams.width = 512;
    snpParams.height = 512;
    snpParams.ppi_ratio = 1;
    snpParams.lat = 0;
    snpParams.lng = 0;
    snpParams.zoom = 0;
    snpParams.pitch = 0;
    snpParams.bearing = 0;

    struct snapshot_result ret;
    ret = snapshot(snpParams);
    if(ret.did_error){
	    std::cout << "Error :" << ret.err << std::endl;
	    exit(1);
    }
    try {
        std::ofstream out(output, std::ios::binary);
        out << encodePNG(ret.image);
        out.close();
    }catch(std::exception& e) {
	    std::cout << "Error :" << e.what() << std::endl;
	    exit(1);
    }
    return 0;
}
