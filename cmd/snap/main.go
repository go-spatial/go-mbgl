package main

import (
	"fmt"
	"os"

	"github.com/go-spatial/go-mbgl/cmd/snap/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
