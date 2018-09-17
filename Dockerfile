FROM golang:latest

RUN apt-get update
RUN apt-get install -y curl git build-essential zlib1g-dev automake \
                     libtool xutils-dev make cmake pkg-config python-pip \
                     libcurl4-openssl-dev libpng-dev libsqlite3-dev \
                     libllvm3.9

RUN apt-get install -y cmake cmake-data
RUN apt-get install -y ccache

RUN apt-get install -y libxi-dev libglu1-mesa-dev x11proto-randr-dev \
                     x11proto-xext-dev libxrandr-dev \
                     x11proto-xf86vidmode-dev libxxf86vm-dev \
                     libxcursor-dev libxinerama-dev

RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get install -y nodejs

RUN go get gihub.com/go-spatial/go-mbgl

# RUN mkdir -p /go/src/github.com/go-spatial/go-mbgl
WORKDIR  /go/src/github.com/go-spatial/go-mbgl

ENTRYPOINT bash