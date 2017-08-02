package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configureListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return ca-cli profile list",
	Long:  `Use the command to get ca-cli profile list with local config file.`,
	RunE:  execConfigureList,
}

func init() {
	configureCmd.AddCommand(configureListCmd)
}

func execConfigureList(cmd *cobra.Command, args []string) error {
	// get profile data
	currentConfig, err := getConfig()
	if err != nil {
		return err
	}
	fmt.Println("\nCurrent \"config\" value list:")
	drawConfigTable("", currentConfig)

	return nil
}
