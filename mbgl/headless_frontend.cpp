#include <iostream>

#include <mbgl/gl/headless_frontend.hpp>
#include <mbgl/storage/file_source.hpp>
#include <mbgl/actor/scheduler.hpp>
#include <mbgl/map/map.hpp>

#include "headless_frontend.h"
#include "map.h"
#include "size.h"
#include "file_source.h"
#include "scheduler.h"

using namespace mbgl;


MbglHeadlessFrontend * mbgl_headless_frontend_new(
	MbglSize * size,
	float pixelRatio,
	MbglFileSource * source,
	MbglScheduler * sched,
	const char * cacheDir,
	const char * fontFamily) {

	auto _size = reinterpret_cast<Size *>(size);
	auto _src = reinterpret_cast<FileSource *>(source);
	auto _sched = reinterpret_cast<Scheduler *>(sched);

	optional<std::string> _cacheDir;
	if (cacheDir != nullptr) {
		std::cout << "cachedir not null" << std::endl;
		_cacheDir = std::string(cacheDir);
	}

	optional<std::string> _fontFamily;
	if (fontFamily != nullptr) {
		std::cout << "fontfamily not null" << std::endl;
		_fontFamily = std::string(fontFamily);
	}

	std::cout << "new headless frontend" << std::endl;

	auto front = new mbgl::HeadlessFrontend( *_size, pixelRatio, *_src, *_sched);
	std::cout << "got headless frontend:" << front << std::endl;

	return reinterpret_cast<MbglHeadlessFrontend*>(front);
}

void mbgl_headless_frontend_destruct(MbglHeadlessFrontend * self) {
	std::cout << "Destructor" << self << std::endl;
	auto _self = reinterpret_cast<mbgl::HeadlessFrontend *>(self);
	std::cout << "Before delete" << _self << std::endl;
	delete _self;
	std::cout << "after delete" << std::endl;
}

MbglPremultipliedImage * mbgl_headless_frontend_render(
	MbglHeadlessFrontend * self,
	MbglMap * map) {
	auto _self = reinterpret_cast<mbgl::HeadlessFrontend *>(self);
	auto _map = reinterpret_cast<mbgl::Map *>(map);

	std::cout << "calling _self->render" << std::endl;
	auto img = _self->render(*_map);
	std::cout << "self rendered" << std::endl;
	auto img2 = new mbgl::PremultipliedImage();
	*img2 = std::move(img);

	return reinterpret_cast<MbglPremultipliedImage *>(img2);
}
