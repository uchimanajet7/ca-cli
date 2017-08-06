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
var awsAccountDeleteID string

var awsAccountDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete CA AWS account",
	Long:  `Use the "/aws_accounts/{id}" API to delete specified ID AWS account with Cloud Automator.`,
	RunE:  execAwsAccountDelete,
}

func init() {
	awsAccountCmd.AddCommand(awsAccountDeleteCmd)
	awsAccountDeleteCmd.Flags().StringVar(&awsAccountDeleteID, "aws-id", "", "Specify CA AWS account ID. [required]")
}

func execAwsAccountDelete(cmd *cobra.Command, args []string) error {
	paramID := awsAccountDeleteID
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
	method := "DELETE"

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
