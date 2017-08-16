package cmd

import "github.com/spf13/cobra"

var postProcessCmd = &cobra.Command{
	Use:   "post-process",
	Short: "Manage CA psot processes",
	Long:  `Use the "/post_processes" API to manage psot processes with Cloud Automator.`,
}

func init() {
	RootCmd.AddCommand(postProcessCmd)
}
