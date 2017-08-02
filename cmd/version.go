package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ca-cli",
	Long:  `All software has versions. This is ca-cli's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud Automator CLI v0.1 -- HEAD")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
