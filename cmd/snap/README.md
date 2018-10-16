# Snap
Snap aims to be a headless raster tile, and static map server. Currenly using OSMesa 6, with CPU only rendering.


## URLs that are supported by the application.



### Raster Tile Server API

```
/styles/${style-name:string}/tiles/[${tilesize:int}/]${z:int}/${x:int}/${y:int}[@2x][.${file-extention:enum(jpg,png)}]
```

* style-name     [required] : the name of the style. If loaded via the command line the style name will be "default" (currently this 
the only thing that is supported.) 
* tilesize       [optional] : Default is 512, valid valus are positive multiples of 256. 
* z              [required] : the zoom
* x              [required] : the x coordinate (column) in the slippy tile scheme.
* y              [required] : the y coordinate (row) in the slippy tile scheme.
* @2x            [optional] : to serve hight definition (retina) tiles. Omit to serve standard definition tiles.
* file-extension [optional] : the file type to encode the raster image in. Currently supported formats png, jpg. Default is jpg.

### Static map server

For generating an image of a map at a given point and zoom use the following url.

```
/styles/${style-name:string}/static/${lon:float},${lat:float},${zoom:float},[${bearing:float],[${pitch:float}]]/${width:int}x${height:int}[@2x][.${file-extention:enum(jpg,png)}]
```

* style-name     [required] : the name of the style. If loaded via the command line the style name will be "default" (currently this 
the only thing that is supported.) 
* lon       [optional] : Default is 512, valid valus are positive multiples of 256. 
* lat       [required] : the zoom
* zoom      [required] : the x coordinate (column) in the slippy tile scheme.
* bearing   [required] : the y coordinate (row) in the slippy tile scheme.
* @2x            [optional] : to serve hight definition (retina) tiles. Omit to serve standard definition tiles.
* file-extension [optional] : the file type to encode the raster image in. Currently supported formats are png, jpg. Default is jpg.

### Health check

If the server is up this url will return a 200.

```
/health
```

To run `snap` you can use the following subcomamnds

* serve to run the raster tile server
* generate to generate an image 

## How to build.

Currently snap only supports Linux, and has only been tested on Ubuntu 18.04.

`go build` in the snap directory will build you a new binary using OSMesa. 

### Library dependencies:

This utility depends on `libosmesa6` library.
For Ubuntu 18.04 it can be installed with the following command: `apt-get install libosmesa6`

### The utility uses prebuilt libraries

These libraries are stored in the `mbgl/mbgl/c/lib/linux` directory. To rebuild build
these libraries use the install.sh script in the `mbgl/mbgl/c/` directory.


