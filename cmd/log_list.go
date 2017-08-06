package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "link" flag
var logListLink string

var logListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return CA job log list",
	Long:  `Use the "/logs" API to get job log list with Cloud Automator.`,
	RunE:  execLogList,
}

func init() {
	logCmd.AddCommand(logListCmd)
	logListCmd.Flags().StringVar(&logListLink, "link", "", "Specify CA paging URL of API response.")
}

func execLogList(cmd *cobra.Command, args []string) error {
	paramlink := logListLink

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := "logs"
	method := "GET"
	if paramlink != "" {
		endURL = paramlink
		apiPath = ""
	}

	// create api client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/logs]
	httpRequest, err := client.createRequest(ctx, method, apiPath, currentProfile.APIKey, nil)
	if err != nil {
		return err
	}

	httpResponse, err := client.HTTPClient.Do(httpRequest)
	if err != nil {
		return err
	}

	// get result string
	result, err := getBodyString(httpResponse)
	if err != nil {
		return err
	}

	// draw result
	fmt.Printf("%s\n", result)

	return nil
}
