package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "link" flag
var awsAccountListLink string

var awsAccountListCmd = &cobra.Command{
	Use:   "list",
	Short: "Return CA AWS account list",
	Long:  `Use the "/aws_accounts" API to get AWS account list with Cloud Automator.`,
	RunE:  execAwsAccountList,
}

func init() {
	awsAccountCmd.AddCommand(awsAccountListCmd)
	awsAccountListCmd.Flags().StringVar(&awsAccountListLink, "link", "", "Specify CA paging URL of API response.")
}

func execAwsAccountList(cmd *cobra.Command, args []string) error {
	paramlink := awsAccountListLink

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := "aws_accounts"
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

	// call API [/aws_accounts]
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
