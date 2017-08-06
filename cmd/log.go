package cmd

import "github.com/spf13/cobra"

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Manage CA job logs",
	Long:  `Use the "/logs" API to manage job logs with Cloud Automator.`,
}

func init() {
	RootCmd.AddCommand(logCmd)
}
