package cmd

import (
	"fmt"
	"regexp"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Manage ca-cli profiles",
	Long:  `Use the command to manage and register ca-cli profiles with local config file.`,
	RunE:  execConfigure,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func execConfigure(cmd *cobra.Command, args []string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	paramProfileName := rootPprofileName
	fmt.Printf("Register the information necessary for execution as a profile of \"%s\".\n\n", paramProfileName)

	// get config data
	currentConfig, err := getConfig()
	if err == nil {
		// get profile data
		_, err := getProfile(paramProfileName)
		if err == nil {
			// if there is existing value
			fmt.Printf("Current \"%s\" profile value: \n", paramProfileName)
			drawConfigTable(paramProfileName, currentConfig)

			if !prompter.YN(fmt.Sprintf("Overwrite current \"%s\" profile value?", paramProfileName), false) {
				return nil
			}
		}
	}

	// no echo + 32 characters or more
	apiKey := (&prompter.Prompter{
		Message: "API Key",
		Regexp:  regexp.MustCompile(`.{32,}`),
		NoEcho:  true,
	}).Prompt()
	fmt.Printf("API Key: %s\n", maskedAPIKey(apiKey))

	// input custom endpoint
	endpoint := prompter.Prompt("Endpoint", "")

	// current config not exist
	if err != nil {
		currentConfig = &config{
			Endpoint: getEndpoint(),
			Profiles: map[string]profile{
				paramProfileName: profile{
					APIKey:   apiKey,
					Endpoint: endpoint,
				},
			},
		}
	} else {
		// add current config
		currentConfig.Profiles[paramProfileName] = profile{
			APIKey:   apiKey,
			Endpoint: endpoint,
		}
	}

	return saveConfig(path, currentConfig)
}
