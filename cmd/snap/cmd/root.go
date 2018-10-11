package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "snap",
	Short: "snap is a raster tile server",
	Long:  `snap is a raster tile server version: ` + Version,
}

const DefaultRootStyle = "https://raw.githubusercontent.com/go-spatial/tegola-web-demo/master/styles/hot-osm.json"

var RootStyle string = DefaultRootStyle

func init() {

	RootCmd.PersistentFlags().StringVarP(&RootStyle, "style", "s", DefaultRootStyle, "style to use. Style name will be default")

	RootCmd.AddCommand(cmdServer)
	RootCmd.AddCommand(cmdVersion)
	RootCmd.AddCommand(cmdGenerate)
}
