package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "version not set"

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number.",
	Long:  `The version of the software. [` + Version + `]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
