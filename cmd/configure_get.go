package cmd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/spf13/cobra"
)

var configureGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Return ca-cli profile",
	Long:  `Use the command to get specified name ca-cli profile with local config file.`,
	RunE:  execConfigureGet,
}

func init() {
	configureCmd.AddCommand(configureGetCmd)
}

func execConfigureGet(cmd *cobra.Command, args []string) error {
	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// get struct field names
	rt := reflect.TypeOf(*currentProfile)
	structNames := make([]string, 0, rt.NumField()+1)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		structNames = append(structNames, field.Name)
	}

	input := (&prompter.Prompter{
		Choices:    structNames,
		Default:    structNames[0],
		Message:    "Please enter what you want to display",
		IgnoreCase: true,
	}).Prompt()

	var drawText string
	switch strings.ToLower(input) {
	case "apikey":
		drawText = fmt.Sprintf("APIKey: %s", currentProfile.APIKey)
	case "endpoint":
		drawText = fmt.Sprintf("Endpoint: %s", currentProfile.Endpoint)
	}

	fmt.Printf("\nCurrent \"%s\" profile value: \n", paramProfileName)
	fmt.Println(drawText)

	return nil
}
