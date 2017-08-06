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

// holds value of "aws-id" flag
var awsAccountGetID string

var awsAccountGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Return CA AWS account",
	Long:  `Use the "/aws_accounts/{id}" API to get specified ID AWS account with Cloud Automator.`,
	RunE:  execAwsAccountGet,
}

func init() {
	awsAccountCmd.AddCommand(awsAccountGetCmd)
	awsAccountGetCmd.Flags().StringVar(&awsAccountGetID, "aws-id", "", "Specify CA AWS account ID. [required]")
}

func execAwsAccountGet(cmd *cobra.Command, args []string) error {
	paramID := awsAccountGetID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--aws-id] required CA AWS account ID.")
	}

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := path.Join("aws_accounts", paramID)
	method := "GET"

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/aws_accounts/{id}]
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
