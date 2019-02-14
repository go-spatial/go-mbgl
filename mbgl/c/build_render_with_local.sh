#!/bin/bash -e

echo "Building the object file."

g++ \
	-DRAPIDJSON_HAS_STDSTRING=1 \
	-D_GLIBCXX_USE_CXX11_ABI=1 \
	-Iinclude \
	-Iinclude/include \
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
	lib/linux/libmbgl-core.a \
	lib/linux/libmbgl-filesource.a \
	lib/linux/libmbgl-loop-uv.a \
	lib/linux/libsqlite.a \
	-Wl,--no-as-needed -lcurl \
	-Wl,--as-needed lib/linux/libmbgl-core.a \
	-lOSMesa \
	lib/linux/libpng.a \
	-lz \
	-lm \
	lib/linux/libjpeg.a \
	lib/linux/libnunicode.a \
	lib/linux/libicu.a \
	-Wl,-Bsymbolic-functions lib/linux/libuv.a \
	-lrt \
	-lpthread \
	-lnsl \
	-ldl \
	-o bin/mbgl-render-stripped \

