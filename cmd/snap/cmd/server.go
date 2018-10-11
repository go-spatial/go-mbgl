package cmd

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimfeld/httptreemux"
	"github.com/disintegration/imaging"
	"github.com/go-spatial/geom/slippy"

	"github.com/go-spatial/go-mbgl/internal/bounds"
	mbgl "github.com/go-spatial/go-mbgl/mbgl/simplified"

	"github.com/spf13/cobra"
)

var cmdServer = &cobra.Command{
	Use:     "serve",
	Short:   "Use snap as a tile server",
	Aliases: []string{"server"},
	Long: `Use snap as a tile server. Maps tiles will be served at following urls:
	 /styles/:style-name/tiles/[tilesize]/:z/:x/:y[@2x].[file-extension]
	 or
	 /styles/:style-name/static/:lon,:lat,:zoom[,:bearing[,:pitch]]/:widthx:height[@2x][.:file-extension]
	 where 
	 • style-name [required]: the name of the style. If loaded via the command line
	 the style name will be "default." If loaded via a config file the name of the
	 style to reference -- "default" for a config will be the first style in the
	 config file.
	 • tilesize [optional]: Default is 512x512 pixels.
	 • z [required]: the zoom
	 • x [required]: the x coordinate (column) in the slippy tile scheme.
	 • y [required]: the y coordinate (row) in the slippy tile scheme.
	 • @2 [optional]: to server high defination (retina) tiles. Omit to serve standard definition tiles.
	 • file-extension [optional]: the file type to encode the raster image in. Values: png, jpg. Default jpg.
	 • lon [required]: Longitude for the center point of the static map. -180 and 180.
	 • lat [required]: Latitude for the center point of the static map. -90 and 90.
	 • zoom [required]: Zoom level; a number between 0 and 20.
	 • bearing [optional]: Bearing rotates the map around it center. An number between 0 and 360; default 0.
	 • pitch [optional]: Pitch tilts the map, producing a perspective effect. A number between 0 and 60; default 0.
	 • width [required]: Width of the image; a number between 1 and 1280 pixels.
	 • height [required]: Height of the image; a number between 1 and 1280 pixels.
`,
	Run: commandServer,
}

const defaultServerAddress = ":8080"

var cmdServerAddress string = defaultServerAddress

func init() {
	cmdServer.Flags().StringVarP(&cmdServerAddress, "address", "a", defaultServerAddress, "address to bind the tile server to.")
}

func commandServer(cmd *cobra.Command, args []string) {
	fmt.Println("Would start up the server here.", strings.Join(args, " , "))
	// start our server
	router := newRouter()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !strings.Contains(RootStyle, "://") {
		RootStyle = "file://" + RootStyle
	}

	mbgl.StartSnapshotManager(ctx)

	http.ListenAndServe(cmdServerAddress, router)
}

func newRouter() *httptreemux.TreeMux {

	r := httptreemux.New()
	r.GET("/health", serverHandlerHealth)

	group := r.NewGroup("/styles")
	group.GET("/:style-name/tiles/:tilesize/:z/:x/:ypath", serverHandlerTileSize)
	group.GET("/:style-name/tiles/:z/:x/:ypath", serverHandlerTile)
	group.GET("/:style-name/static/:lonlatzoompath/:widthheightpath", serverHandlerStatic)

	return r
}

func parsePPIPath(path []byte) (startPos int, val float64, err error) {

	// look for the last @...x
	startPos = bytes.LastIndexByte(path, '@')
	if startPos == -1 {
		return -1, 1, nil
	}

	endPos := bytes.IndexByte(path[startPos:], 'x')
	if endPos == -1 {
		return -1, 1, nil
	}

	val, err = strconv.ParseFloat(string(path[startPos+1:startPos+endPos]), 64)
	return startPos, val, err
}

func parseAtAndDot(path []byte) (pre []byte, ppi float64, ext string, err error) {

	ext = "jpg"

	extPos := bytes.LastIndexByte(path, '.')
	if extPos != -1 && extPos+1 < len(path) {
		ext = string(bytes.ToLower(path[extPos+1:]))
	}

	atPos, ppi, err := parsePPIPath(path)
	if err != nil {
		return pre, 1.0, ext, err
	}

	switch {
	case atPos != -1:
		pre = path[:atPos]
	case extPos != -1:
		pre = path[:extPos]
	default:
		pre = path
	}

	return pre, ppi, ext, nil
}

func parseWidthHeightPath(widthHeightPath string) (width, height int, at2x float64, ext string, err error) {

	prePath, at2x, ext, err := parseAtAndDot([]byte(widthHeightPath))
	if err != nil {
		return 0, 0, at2x, ext, fmt.Errorf("Error parsing extention or ppi (%v): %v", widthHeightPath, err)
	}

	// width x height
	xindex := bytes.IndexByte(prePath, 'x')
	if xindex == -1 {
		return 0, 0, at2x, ext, fmt.Errorf("expected width 'x' height. Did not find seperator. %v", string(prePath))
	}

	w64, err := strconv.ParseInt(string(prePath[:xindex]), 10, 64)
	if err != nil {
		return 0, 0, at2x, ext, fmt.Errorf("Error parsing width(%v): %v", string(prePath[:xindex]), err)
	}

	h64, err := strconv.ParseInt(string(prePath[xindex+1:]), 10, 64)
	if err != nil {
		return 0, 0, at2x, ext, fmt.Errorf("Error parsing height(%v): %v", string(prePath[xindex+1:]), err)
	}

	return int(w64), int(h64), at2x, ext, nil

}

func parseYPath(ypath string) (y uint, at2x float64, ext string, err error) {

	prePath, at2x, ext, err := parseAtAndDot([]byte(ypath))
	if err != nil {
		return y, at2x, ext, fmt.Errorf("Error parsing extention or ppi (%v): %v", ypath, err)
	}
	y64, err := strconv.ParseUint(string(prePath), 10, 64)
	return uint(y64), at2x, ext, err

}

func centerZoom(tilesize int, z, x, y uint) (center [2]float64, zoom float64) {

	var tile = slippy.Tile{
		Z: z,
		X: x,
		Y: y,
	}
	center = bounds.Center(tile.Extent4326())
	n := int(math.Log2(float64(tilesize / 256)))
	zoom = float64(int(z) - n)
	if zoom < 0.0 {
		zoom = 0.0
	}
	return center, zoom
}

func serverHandlerHealth(w http.ResponseWriter, r *http.Request, params map[string]string) {
	w.WriteHeader(http.StatusOK)
}

// serverhandlerTileSize will handle the tiles url with a tilesize value in it.
func serverHandlerTileSize(w http.ResponseWriter, r *http.Request, params map[string]string) {

	y, ppi, ext, err := parseYPath(params["ypath"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	x64, err := strconv.ParseUint(strings.TrimSpace(params["x"]), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse param x: %v", err.Error()), http.StatusBadRequest)
		return
	}
	if x64 < 0 {
		http.Error(w, fmt.Sprintf("Failed param x should be greater then 0: %v", x64), http.StatusBadRequest)
		return
	}
	z64, err := strconv.ParseUint(strings.TrimSpace(params["z"]), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse param z: %v", err.Error()), http.StatusBadRequest)
		return
	}
	if z64 < 0 {
		http.Error(w, fmt.Sprintf("Failed param z should be greater then 0: %v", z64), http.StatusBadRequest)
		return
	}
	tileSize64, err := strconv.ParseUint(strings.TrimSpace(params["tilesize"]), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse param tile size: %v", err.Error()), http.StatusBadRequest)
		return
	}
	if tileSize64%256 != 0 {
		http.Error(w, fmt.Sprintf("Failed param tile-size should be a multiple of 256: %v", tileSize64), http.StatusBadRequest)
		return
	}
	styleName := strings.ToLower(strings.TrimSpace(params["style-name"]))
	if styleName != "default" {
		http.Error(w, fmt.Sprintf("Unknown style %v", styleName), http.StatusNotFound)
		return
	}

	if ext != "jpg" && ext != "png" {
		http.Error(w, fmt.Sprintf("only supported extentions are jpg and png, got [%v]\n", ext), http.StatusBadRequest)
		return
	}

	buffer := new(bytes.Buffer)
	imgType, err := generateTileImage(buffer, RootStyle, float64(tileSize64), uint(z64), uint(x64), y, ppi, ext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate %v image: %v", ext, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", imgType)
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
	if _, err = io.Copy(w, buffer); err != nil {
		log.Println("Got error writing to write out image:", err)
	}

}

func generateTileImage(w io.Writer, style string, tsize float64, z, x, y uint, ppi float64, ext string) (string, error) {

	center, zoom := centerZoom(int(tsize), z, x, y)

	tilesize := uint32(float64(tsize) * ppi)
	// .20 from expermintation
	var pxBuffer uint32
	switch {
	case zoom >= 2:
		pxBuffer = uint32(float64(tsize) * 0.10)
	case zoom >= 10:
		pxBuffer = uint32(float64(tsize) * 0.20)
	}

	snpsht := mbgl.Snapshotter{
		Style:    style,
		Width:    tilesize + pxBuffer,
		Height:   tilesize + pxBuffer,
		PPIRatio: ppi,
		Lat:      center[0],
		Lng:      center[1],
		Zoom:     zoom,
	}

	return generateImage(w, snpsht, int(tilesize), int(tilesize), ext)
}

func generateZoomCenterImage(w io.Writer, style string, width, height int, ppi, pitch, bearing float64, center [2]float64, zoom float64, ext string) (string, error) {
	_width := uint32(float64(width) * ppi)
	_height := uint32(float64(height) * ppi)
	snpsht := mbgl.Snapshotter{
		Style:    style,
		Width:    _width,
		Height:   _height,
		PPIRatio: ppi,
		Lat:      center[1],
		Lng:      center[0],
		Zoom:     zoom,
		Pitch:    pitch,
		Bearing:  bearing,
	}

	return generateImage(w, snpsht, int(_width), int(_height), ext)
}

func generateImage(w io.Writer, param mbgl.Snapshotter, width, height int, ext string) (string, error) {
	imageType := "unknown"
	param.Zoom -= 1

	img, err := mbgl.Snapshot1(param)
	if err != nil {
		return imageType, err
	}

	cimg := imaging.CropCenter(img, width, height)

	switch ext {
	case "png":
		imageType = "image/png"
		if err := png.Encode(w, cimg); err != nil {
			return imageType, err
		}
	case "jpg":
		imageType = "image/jpeg"
		if err := jpeg.Encode(w, cimg, nil); err != nil {
			return imageType, err
		}
	}
	return imageType, nil
}

func serverHandlerTile(w http.ResponseWriter, r *http.Request, params map[string]string) {
	params["tilesize"] = "256"
	serverHandlerTileSize(w, r, params)
}

func serverHandlerStatic(w http.ResponseWriter, r *http.Request, params map[string]string) {
	// /styles/:style-name/static/:lon,:lat,:zoom[,:bearing[,:pitch]]/:widthx:height[@2x][.:file-extension]

	styleName := strings.ToLower(strings.TrimSpace(params["style-name"]))
	if styleName != "default" {
		http.Error(w, fmt.Sprintf("Unknown style %v", styleName), http.StatusNotFound)
		return
	}
	lnglatzParts := strings.Split(params["lonlatzoompath"], ",")
	if len(lnglatzParts) < 3 {
		http.Error(w, fmt.Sprintf("not enough params for lat lng and zoom: got %v", params["lonlatzoompath"]), http.StatusBadRequest)
		return
	}
	lng64, err := strconv.ParseFloat(lnglatzParts[0], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse lng %v", err), http.StatusBadRequest)
		return
	}
	lat64, err := strconv.ParseFloat(lnglatzParts[1], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse lat %v", err), http.StatusBadRequest)
		return
	}
	zoom, err := strconv.ParseFloat(lnglatzParts[2], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse zoom %v", err), http.StatusBadRequest)
		return
	}
	var bear64, pitch64 float64
	if len(lnglatzParts) >= 4 {
		bear64, err = strconv.ParseFloat(lnglatzParts[3], 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse bearing %v", err), http.StatusBadRequest)
			return
		}
	}
	if len(lnglatzParts) >= 5 {
		pitch64, err = strconv.ParseFloat(lnglatzParts[4], 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not parse pitch %v", err), http.StatusBadRequest)
			return
		}
	}
	width, height, ppi, ext, err := parseWidthHeightPath(params["widthheightpath"])
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse width height %v", err), http.StatusBadRequest)
		return
	}

	if ext != "jpg" && ext != "png" {
		http.Error(w, fmt.Sprintf("only supported extentions are jpg and png, got [%v]\n", ext), http.StatusBadRequest)
		return
	}

	width = int(float64(width) * ppi)
	height = int(float64(height) * ppi)

	snpsht := mbgl.Snapshotter{
		Style:    RootStyle,
		Width:    uint32(width),
		Height:   uint32(height),
		PPIRatio: ppi,
		Lat:      lat64,
		Lng:      lng64,
		Zoom:     zoom,
		Pitch:    pitch64,
		Bearing:  bear64,
	}

	buffer := new(bytes.Buffer)

	imgType, err := generateImage(buffer, snpsht, width, height, ext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate %v image: %v", ext, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", imgType)
	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
	if _, err = io.Copy(w, buffer); err != nil {
		log.Println("Got error writing to write out image:", err)
	}

}
