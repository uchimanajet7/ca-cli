package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version number is set automatically at build time
var version string

// revision is set automatically use "git rev-parse --short HEAD" at build time
var revision string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ca-cli",
	Long:  `All software has versions. This is ca-cli's`,
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "[Sorry none set at build time.]"
		}
		if revision == "" {
			revision = "[Sorry none set at build time.]"
		}

		fmt.Printf("Cloud Automator CLI %s (-- HEAD rev: %s)\n", version, revision)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
