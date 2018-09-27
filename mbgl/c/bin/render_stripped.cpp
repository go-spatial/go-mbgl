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

int main(int argc __attribute__((unused)), char *argv[] __attribute__((unused)) ) {
    const std::string output     = "out.png";
    const std::string cache_file = "cache.sqlite";
    const std::string asset_root = ".";

    using namespace mbgl;

    util::RunLoop loop;

    DefaultFileSource fileSource(cache_file, asset_root);
    ThreadPool threadPool(4);
    HeadlessFrontend frontend({ 512, 512 }, 1, fileSource, threadPool);
    Map map(frontend, MapObserver::nullObserver(), frontend.getSize(), 1, fileSource, threadPool, MapMode::Static);


    map.getStyle().loadURL("file://style.json");
    map.setLatLngZoom({ 0, 0 }, 0);
    map.setBearing(0);
    map.setPitch(0);


    try {
        std::ofstream out(output, std::ios::binary);
        out << encodePNG(frontend.render(map));
        out.close();
    } catch(std::exception& e) {
        std::cout << "Error: " << e.what() << std::endl;
        exit(1);
    }

    return 0;
}
