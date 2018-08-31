#!/usr/bin/env bash

# check mbgl dependencies
function check_dep {
    arg=$1
    which  $arg #> /dev/null
    if [[ ! $? -eq 0 ]]; then
        echo "dep $arg, not installed"
        exit
    fi
}

deps="node cmake ccache xcpretty jazzy"
for dep in $deps; do
    echo "checking if $dep is installed"
    check_dep $dep
done

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

mkdir $PKG_ROOT/lib
cd $PKG_ROOT/mapbox-gl-native

if [[ uname -eq Darwin ]]; then
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
    echo "no install instructions for $(uname)"
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