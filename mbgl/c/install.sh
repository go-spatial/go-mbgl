#!/bin/bash -x


# copy includes will use find and xargs to copy includes from the mason_packages/headers folder
# to the $PKG_ROOT/include directory
# find "mapbox-gl-native/mason_packages/headers/geometry" -iname include | xargs -n 1 -I {} echo "DIR id {} stuff"
function copy_includes {
	subdir=$1
	shift
	for libraryName in $@; do
	  srcdir="${PKG_ROOT}/mapbox-gl-native/${subdir}/${libName}"
	  echo "copying includes for ${PKG_ROOT}/mapbox-gl-native/${subdir}/${libraryName}"
 	  find "${PKG_ROOT}/mapbox-gl-native/${subdir}/${libraryName}" -iname include | xargs -n 1 -I{} cp -R {}/ ${PKG_ROOT}/
	done	
}

function copy_hpps {
	subdir=$1
	shift
	echo "remainder $@"
	for dir in $@; do
	  srcdir="${PKG_ROOT}/mapbox-gl-native/${subdir}/${dir}"
	  for file in $(find ${srcdir} -type f -name '*.hpp'); do
		destfile=${file#$srcdir}
		bdir=$(dirname ${destfile})
		echo copying $(basename $file) " to ${INCLUDEDIR}/$destfile"
		mkdir -p ${INCLUDEDIR}/${bdir}
		cp $file ${INCLUDEDIR}/${destfile}
	   done
	done	
}

function copy_libs {
      subdir=$1
      shift
      for libName in $@; do 
	srcdir="${PKG_ROOT}/mapbox-gl-native/${subdir}/${osdir}/${libName}"
	for afile in $( find "${srcdir}" -iname "*.a" ); do
		cp -R ${afile} ${LIBDIR}/
	done
      done
}

unamestr=`uname`
echo "running on $(uname)"

# check mbgl dependencies
check="which"
if [[ $unamestr == "Linux" ]]; then
    check="dpkg -l"
    echo "running linux, checking deps with $check"
    osdir="linux-x86_64"
fi

function check_dep {
    arg=$1
    $check  $arg #> /dev/null
    err=$?
    if [[ ! err -eq 0 ]]; then
        echo "missing dependency $arg ($err)"
        exit
    fi
}

function no_call1 {
if [[ $unamestr == "Darwin" ]]; then
    deps="node cmake ccache xcpretty jazzy"
    for dep in $deps; do
        echo "checking if $dep is installed"
        check_dep $dep
    done
else
    deps="curl git build-essential zlib1g-dev automake
    libtool xutils-dev make cmake pkg-config python-pip
    libcurl4-openssl-dev 
    libllvm3.9
    libxi-dev libglu1-mesa-dev x11proto-randr-dev
    x11proto-xext-dev libxrandr-dev
    x11proto-xf86vidmode-dev libxxf86vm-dev
    libxcursor-dev libxinerama-dev nodejs"

    # libpng-dev libsqlite3-dev 

    for dep in $deps; do
        echo "checking if $dep is installed"
        check_dep $dep
    done
fi
}

if [[ ! $GOPATH ]]; then
    echo "GOPATH must be set"
    exit
fi

PKG_ROOT=$GOPATH/src/github.com/go-spatial/go-mbgl/mbgl/c


# download and install sdk
if [[ ! -d mapbox-gl-native ]]; then
    git clone --depth=1 --recursive https://github.com/mapbox/mapbox-gl-native
	 pushd mapbox-gl-native
	 git config user.email "gautam.dey77@gmail.com"
	 git config user.name "Gautam Dey"
	 git am ../patches/*
	 popd

fi

if [[ ! -d $PKG_ROOT/lib ]]; then
    mkdir $PKG_ROOT/lib
fi

cd $PKG_ROOT/mapbox-gl-native

if [[ uname == "Darwin" ]]; then
    echo "installing for $(uname)"
    LIBDIR=$PKG_ROOT/lib/darwin

    git checkout tags/macos-v0.9.0
    git submodule foreach --recursive git --reset --hard

    make xpackage
    err=$?
    if [[ ! err -eq 0 ]]; then
        echo "error $err"
        exit
    fi

    if [[ -d ${LIBDIR} ]]; then
        rm -rf ${LIBDIR}
    fi

    mkdir -p ${LIBDIR}

    mv $PKG_ROOT/mapbox-gl-native/build/macos/Debug/* ${LIBDIR}
    sudo mv $PKG_ROOT/lib/darwin/Mapbox.framework /Library/Frameworks/
else
    LIBDIR=$PKG_ROOT/lib/linux
    INCLUDEDIR=$PKG_ROOT/include

    #git checkout master
    #git reset --hard --recurse-submodules
	 
    make WITH_OSMESA=ON linux-core

    if [[ -d ${LIBDIR} ]]; then
        rm -rf ${LIBDIR}
	rm -rf ${INCLUDEDIR}
    fi

    mkdir -p ${LIBDIR}
    mkdir -p ${INCLUDEDIR}

    copy_libs "build" "Debug"
    copy_libs "mason_packages" "libuv" "libjpeg-turbo" "libpng"

    #cp $PKG_ROOT/mapbox-gl-native/build/linux-x86_64/Debug/*.a ${LIBDIR}
    cp -R $PKG_ROOT/mapbox-gl-native/include/* ${INCLUDEDIR}

    copy_hpps "platform" "default"
    
    copy_includes "vendor" "expected"
    copy_includes "mason_packages/headers" "geometry" "variant"
fi

# install mason-js (mapbox package manager)
#cd $PKG_ROOT
#git clone https://github.com/mapbox/mason-js
#cd $PKG_ROOT/mason-js
#npm i -g

# install deps
#cd $PKG_ROOT
#mason-js install
#mason-js link
