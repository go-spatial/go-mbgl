#include <mbgl/style/style.hpp>
#include <mbgl/map/map.hpp>
#include <mbgl/storage/file_source.hpp>
#include <mbgl/actor/scheduler.hpp>

#include "map.h"
#include "size.h"
#include "file_source.h"
#include "scheduler.h"

using namespace mbgl;

MbglMap * mbgl_map_new (MbglRendererFrontend * frontend,
	MbglSize * size,
	float pixelRatio,
	MbglFileSource * src,
	MbglScheduler * sched) {

	auto _frontend = reinterpret_cast<RendererFrontend *>(frontend);
	auto _size = reinterpret_cast<Size *>(size);
	auto _src = reinterpret_cast<FileSource *>(src);
	auto _sched = reinterpret_cast<Scheduler *>(sched);

	auto ptr = new Map(*_frontend, MapObserver::nullObserver(), *_size, pixelRatio, *_src, *_sched, MapMode::Static);

	return reinterpret_cast<MbglMap *>(ptr);
}

void mbgl_map_destruct(MbglMap * self) {
	auto _self = reinterpret_cast<Map *>(self);
	delete _self;
}

void mbgl_map_jump_to(MbglMap * self, MbglCameraOptions * opts) {
	auto _self = reinterpret_cast<Map *>(self);
	auto _opts = reinterpret_cast<CameraOptions *>(opts);

	_self->jumpTo(*_opts);
}


void mbgl_map_set_style_url(MbglMap * self, const char * addr) {
	auto _self = reinterpret_cast<Map *>(self);
	auto _addr = std::string(addr);

	_self->getStyle().loadURL(_addr);
}