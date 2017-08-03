package cmd

import (
	"fmt"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

var configureDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete ca-cli profile",
	Long:  `Use the command to Delete specified name ca-cli profile with local config file.`,
	RunE:  execConfigureDelete,
}

func init() {
	configureCmd.AddCommand(configureDeleteCmd)
}

func execConfigureDelete(cmd *cobra.Command, args []string) error {
	// get profile data
	paramProfileName := rootPprofileName
	_, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// get config data
	currentConfig, err := getConfig()
	if err != nil {
		return err
	}

	// draw current data
	fmt.Printf("\nCurrent \"%s\" profile value: \n", paramProfileName)
	drawConfigTable(paramProfileName, currentConfig)

	// confrim delete
	if !prompter.YN(fmt.Sprintf("Delete current \"%s\" profile?", paramProfileName), false) {
		return nil
	}

	// delete current profile
	delete(currentConfig.Profiles, paramProfileName)

	// save config
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	return saveConfig(path, currentConfig)
}
