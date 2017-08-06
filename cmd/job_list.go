package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "link" flag
var jobListLink string

var jobListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return CA job list",
	Long:  `Use the "/jobs" API to get job list with Cloud Automator.`,
	RunE:  execJobList,
}

func init() {
	jobCmd.AddCommand(jobListCmd)
	jobListCmd.Flags().StringVar(&jobListLink, "link", "", "Specify CA paging URL of API response.")
}

func execJobList(cmd *cobra.Command, args []string) error {
	paramlink := jobListLink

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := "jobs"
	method := "GET"
	if paramlink != "" {
		endURL = paramlink
		apiPath = ""
	}

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/jobs]
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
