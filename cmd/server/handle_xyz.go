package main

import (
	"fmt"
	"image"
	"net/http"
	"strconv"

	"github.com/go-spatial/geom/slippy"
)

func MakeXYZHandler(maps []Map) http.HandleFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handleXYZ(w, r)
	}
}

func handleXYZ(w http.ResponseWriter, r *http.Request) {

	params := LoadParams(r)

	tile := slippy.Tile{}

	arr := []string{"z", "x", "y"}
	for _, v := range arr {
		ui64, err := strconv.ParseUint(params[v], 10, 64)
		if err != nil {
			http.Error(w, "invalid url directory "+params[v], http.StatusBadRequest)
			return
		}

		ui := uint(ui64)
		switch v {
		case "z":
			tile.Z = ui
		case "x":
			tile.X = ui
		case "y":
			tile.Y = ui
		}
	}

	fmt.Println("new request: ", tile.Extent4326(), image.Pt(256, 256))

	img := snap.Snapshot(tile.Extent4326(), image.Pt(256, 256))

	var err error
	switch path.Ext(r.URL.Path) {
	case ".png":
		err = png.Encode(w, img)
	case ".jpg":
		err = jpeg.Encode(w, img, &jpeg.Options{80})
	}

	// if err != nil {
	// 	// todo: log error
	// 	log.Println(err.Error())
	// }
}
