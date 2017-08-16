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

// holds value of "post-process-id" flag
var postProcessDeleteID string

var postProcessDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete CA post process",
	Long:  `Use the "/post_processes/{id}" API to delete specified ID post process with Cloud Automator.`,
	RunE:  execPostProcessDelete,
}

func init() {
	postProcessCmd.AddCommand(postProcessDeleteCmd)
	postProcessDeleteCmd.Flags().StringVar(&postProcessDeleteID, "post-process-id", "", "Specify CA post process ID. [required]")
}

func execPostProcessDelete(cmd *cobra.Command, args []string) error {
	paramID := postProcessDeleteID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--post-process-id] required CA post process ID.")
	}

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := path.Join("post_processes", paramID)
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
