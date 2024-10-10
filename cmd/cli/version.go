package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	buildVersion string
	buildDate    string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Golem",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("Golem " + buildVersion)
	},
}
