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
var postProcessGetID string

var postProcessGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Return CA post process",
	Long:  `Use the "/post_processes/{id}" API to get specified ID post process with Cloud Automator.`,
	RunE:  execPostProcessGet,
}

func init() {
	postProcessCmd.AddCommand(postProcessGetCmd)
	postProcessGetCmd.Flags().StringVar(&postProcessGetID, "post-process-id", "", "Specify CA post process ID. [required]")
}

func execPostProcessGet(cmd *cobra.Command, args []string) error {
	paramID := postProcessGetID
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
	method := "GET"

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
