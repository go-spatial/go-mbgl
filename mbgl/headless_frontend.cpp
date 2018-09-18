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
		_cacheDir = std::string(cacheDir);
	}

	optional<std::string> _fontFamily;
	if (fontFamily != nullptr) {
		_fontFamily = std::string(fontFamily);
	}

	auto front = new HeadlessFrontend(
		*_size,
		pixelRatio,
		*_src,
		*_sched);

	return reinterpret_cast<MbglHeadlessFrontend*>(front);
}

void mbgl_headless_frontend_destruct(MbglHeadlessFrontend * self) {
	auto _self = reinterpret_cast<HeadlessFrontend *>(self);

	delete _self;
}

MbglPremultipliedImage * mbgl_headless_frontend_render(
	MbglHeadlessFrontend * self,
	MbglMap * map) {
	auto _self = reinterpret_cast<HeadlessFrontend *>(self);
	auto _map = reinterpret_cast<Map *>(map);

	auto img = _self->render(*_map);
	auto img2 = new PremultipliedImage();
	*img2 = std::move(img);

	return reinterpret_cast<MbglPremultipliedImage *>(img2);
}