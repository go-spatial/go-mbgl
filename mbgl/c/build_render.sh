#!/bin/bash -e

echo "Building the object file."

g++ \
	-DRAPIDJSON_HAS_STDSTRING=1 \
	-D_GLIBCXX_USE_CXX11_ABI=1 \
	-Imapbox-gl-native/include \
	-Imapbox-gl-native/platform/default \
	-Imapbox-gl-native/mason_packages/headers/boost/1.65.1/include \
	-Imapbox-gl-native/mason_packages/headers/geometry/0.9.3/include \
	-Imapbox-gl-native/mason_packages/headers/variant/1.1.4/include \
	-Imapbox-gl-native/mason_packages/linux-x86_64/libuv/1.9.1/include \
	-Imapbox-gl-native/mason_packages/linux-x86_64/libpng/1.6.25/include/libpng16 \
	-Imapbox-gl-native/mason_packages/linux-x86_64/libjpeg-turbo/1.5.0/include \
	-Imapbox-gl-native/vendor/expected/include \
	-Imapbox-gl-native/vendor/sqlite/include \
	-ftemplate-depth=1024 \
	-Wall \
	-Wextra \
	-Wshadow \
	-Wnon-virtual-dtor \
	-Wno-variadic-macros \
	-Wno-unknown-pragmas \
	-Werror \
	-fext-numeric-literals \
	-g \
	-fPIE \
	-fvisibility=hidden \
	-fvisibility-inlines-hidden \
	-std=c++14 \
	-o bin/render_stripped.cpp.o \
	-c bin/render_stripped.cpp


echo "Building executable."

# libmbgl-core.a libmbgl-filesource.a libmbgl-loop-uv.a libsqlite.a -Wl,--no-as-needed -lcurl -Wl,--as-needed libmbgl-core.a -lOSMesa ../../../mason_packages/linux-x86_64/libpng/1.6.25/lib/libpng.a -lz -lm ../../../mason_packages/linux-x86_64/libjpeg-turbo/1.5.0/lib/libjpeg.a libnunicode.a libicu.a -lz -static-libstdc++ -Wl,-Bsymbolic-functions ../../../mason_packages/linux-x86_64/libuv/1.9.1/lib/libuv.a -lrt -lpthread -lnsl -ldl

g++ \
	-ftemplate-depth=1024 \
	-Wall \
	-Wextra \
	-Wshadow \
	-Wnon-virtual-dtor \
	-Wno-variadic-macros \
	-Wno-unknown-pragmas \
	-Werror -fext-numeric-literals \
	-static-libstdc++ \
	-g bin/render_stripped.cpp.o \
	mapbox-gl-native/build/linux-x86_64/Debug/libmbgl-core.a \
	mapbox-gl-native/build/linux-x86_64/Debug/libmbgl-filesource.a \
	mapbox-gl-native/build/linux-x86_64/Debug/libmbgl-loop-uv.a \
	mapbox-gl-native/build/linux-x86_64/Debug/libsqlite.a \
	-Wl,--no-as-needed -lcurl \
	-Wl,--as-needed mapbox-gl-native/build/linux-x86_64/Debug/libmbgl-core.a \
	-lOSMesa \
	mapbox-gl-native/mason_packages/linux-x86_64/libpng/1.6.25/lib/libpng.a \
	-lz \
	-lm \
	mapbox-gl-native/mason_packages/linux-x86_64/libjpeg-turbo/1.5.0/lib/libjpeg.a \
	mapbox-gl-native/build/linux-x86_64/Debug/libnunicode.a \
	mapbox-gl-native/build/linux-x86_64/Debug/libicu.a \
	-Wl,-Bsymbolic-functions mapbox-gl-native/mason_packages/linux-x86_64/libuv/1.9.1/lib/libuv.a \
	-lrt \
	-lpthread \
	-lnsl \
	-ldl \
	-o bin/mbgl-render-stripped \

