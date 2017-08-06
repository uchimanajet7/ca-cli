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

// holds value of "job-id" flag
var jobDeleteID string

var jobDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete CA job",
	Long:  `Use the "/jobs/{id}" API to delete specified ID job with Cloud Automator.`,
	RunE:  execJobDelete,
}

func init() {
	jobCmd.AddCommand(jobDeleteCmd)
	jobDeleteCmd.Flags().StringVar(&jobDeleteID, "job-id", "", "Specify CA job ID. [required]")
}

func execJobDelete(cmd *cobra.Command, args []string) error {
	paramID := jobDeleteID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--job-id] required CA job ID.")
	}

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := path.Join("jobs", paramID)
	method := "DELETE"

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/jobs/{id}]
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
