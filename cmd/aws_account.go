package cmd

import "github.com/spf13/cobra"

var awsAccountCmd = &cobra.Command{
	Use:   "aws-account",
	Short: "Manage CA AWS accounts",
	Long:  `Use the "/aws_accounts" API to manage AWS accounts with Cloud Automator.`,
}

func init() {
	RootCmd.AddCommand(awsAccountCmd)
}
