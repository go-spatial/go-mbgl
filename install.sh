#!/usr/bin/env bash

unamestr=`uname`
echo "running on $(uname)"

# check mbgl dependencies
check="which"
if [[ $unamestr == "Linux" ]]; then
    check="dpkg -l"
    echo "running linux, checking deps with $check"
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

if [[ $unamestr == "Darwin" ]]; then
    deps="node cmake ccache xcpretty jazzy"
    for dep in $deps; do
        echo "checking if $dep is installed"
        check_dep $dep
    done
else
    deps="curl git build-essential zlib1g-dev automake
    libtool xutils-dev make cmake pkg-config python-pip
    libcurl4-openssl-dev libpng-dev libsqlite3-dev libllvm3.9
    libxi-dev libglu1-mesa-dev x11proto-randr-dev
    x11proto-xext-dev libxrandr-dev
    x11proto-xf86vidmode-dev libxxf86vm-dev
    libxcursor-dev libxinerama-dev nodejs"

    for dep in $deps; do
        echo "checking if $dep is installed"
        check_dep $dep
    done
fi

if [[ ! $GOPATH ]]; then
    echo "GOPATH must be set"
    exit
fi

PKG_ROOT=$GOPATH/src/github.com/go-spatial/go-mbgl


# download and install sdk
if [[ ! -d mapbox-gl-native ]]; then
    git clone https://github.com/mapbox/mapbox-gl-native
fi

git fetch --all --tags --prune

if [[ ! -d $PKG_ROOT/lib ]]; then
    mkdir $PKG_ROOT/lib
fi

cd $PKG_ROOT/mapbox-gl-native

if [[ uname == "Darwin" ]]; then
    echo "installing for $(uname)"

    git checkout tags/macos-v0.9.0
    git reset --hard --recurse-submodules

    make xpackage
    err=$?
    if [[ ! err -eq 0 ]]; then
        echo "error $err"
        exit
    fi

    if [[ -d $PKG_ROOT/lib/darwin ]]; then
        rm -rf $PKG_ROOT/lib/darwin
    fi

    mkdir $PKG_ROOT/lib/darwin

    mv $PKG_ROOT/mapbox-gl-native/build/macos/Debug/* $PKG_ROOT/lib/darwin/
    sudo mv $PKG_ROOT/lib/darwin/Mapbox.framework /Library/Frameworks/
else
    git checkout master
    git reset --hard --recurse-submodules
    make linux-core

    if [[ -d $PKG_ROOT/lib/linux ]]; then
        rm -rf $PKG_ROOT/lib/linux
    fi

    mv $PKG_ROOT/mapbox-gl-native/build/Debug/*.a $PKG_ROOT/lib/linux
fi

# install mason-js (mapbox package manager)
cd $PKG_ROOT
git clone https://github.com/mapbox/mason-js
cd $PKG_ROOT/mason-js
npm i -g

# install deps
cd $PKG_ROOT
mason-js install
mason-js link