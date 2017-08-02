package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// holds value of "profile" flag
var rootPprofileName string

// RootCmd defines 'Cloud Automator(ca)' command
var RootCmd = &cobra.Command{
	Use:           "ca",
	Short:         "Cloud Automator(CA) command",
	Long:          `A command line tool to invoke Cloud Automator(CA) API`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize(initRootConfig)

	RootCmd.PersistentFlags().StringVarP(&rootPprofileName, "profile", "p", getProfileName(), "Specify profile name")
}

func initRootConfig() {
	if _, err := getConfig(); err != nil {
		fmt.Printf("Can't read config: %s\n\n", err)
		// fmt.Print("Please execute the 'configure' command first to register the profile\n\n")
	}
}
