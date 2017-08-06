package cmd

import "github.com/spf13/cobra"

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Manage CA jobs",
	Long:  `Use the "/jobs" API to manage jobs with Cloud Automator.`,
}

func init() {
	RootCmd.AddCommand(jobCmd)
}
