package main

import (
	"flag"
	"net/http"
)

// flags
var (
	style string
)

func main() {
	flag.StringVar(&style, "style", "", "style json which will be used to render")
	flag.Parse()

	mux := Mux{}
	conf := RasterConfig{}

	//snap := mbgl.NewSnapshotterPool(style, 1.0)

	mux.HandleFunc("/raster/maps/:map/:z/:x/:y", MakeXYZHandler(conf.Maps))

	http.ListenAndServe(":8080", mux)
}
