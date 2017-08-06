package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "log-id" flag
var logResourceOperationResultsID string

var logResourceOperationResultsCmd = &cobra.Command{
	Use:   "resource-operation-results",
	Short: "Return CA job log related resource operation results.",
	Long:  `Use the "/logs/{id}/resource_operation_results" API to get specified ID job log related resource operation results with Cloud Automator.`,
	RunE:  execLogResourceOperationResults,
}

func init() {
	logCmd.AddCommand(logResourceOperationResultsCmd)
	logResourceOperationResultsCmd.Flags().StringVar(&logResourceOperationResultsID, "log-id", "", "Specify CA job log ID. [required]")
}

func execLogResourceOperationResults(cmd *cobra.Command, args []string) error {
	paramID := logResourceOperationResultsID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--log-id] required CA job log ID.")
	}

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := path.Join("logs", paramID, "resource_operation_results")
	method := "GET"

	// create api client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/logs/{id}/resource_operation_results]
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
