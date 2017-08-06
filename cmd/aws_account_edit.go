package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "aws-id" flag
var awsAccountEditID string

// holds value of "name" flag
var awsAccountEditName string

// holds value of "account-number" flag
var awsAccountEditAccountNumber string

// holds value of "access-key-id" flag
var awsAccountEditAccessKeyID string

// holds value of "secret-access-key" flag
var awsAccountEditSecretAccessKey string

var awsAccountEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Update CA AWS account",
	Long:  `Use the "/aws_accounts/{id}" API to update specified ID AWS account with Cloud Automator.`,
	RunE:  execAwsAccountEdit,
}

func init() {
	awsAccountCmd.AddCommand(awsAccountEditCmd)

	awsAccountEditCmd.Flags().StringVar(&awsAccountEditID, "aws-id", "", "Specify CA AWS account ID. [required]")
	awsAccountEditCmd.Flags().StringVar(&awsAccountEditName, "name", "", "Specify CA AWS account name.")
	awsAccountEditCmd.Flags().StringVar(&awsAccountEditAccountNumber, "account-number", "", "Specify AWS account number.")
	awsAccountEditCmd.Flags().StringVar(&awsAccountEditAccessKeyID, "access-key-id", "", "Specify AWS access key ID.")
	awsAccountEditCmd.Flags().StringVar(&awsAccountEditSecretAccessKey, "secret-access-key", "", "Specify AWS secret access key.")
}

func execAwsAccountEdit(cmd *cobra.Command, args []string) error {
	paramID := awsAccountEditID
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
	method := "PATCH"

	paramData := map[string]interface{}{}
	paramName := awsAccountEditName
	if paramName != "" {
		paramData["name"] = paramName
	}
	paramNumber := awsAccountEditAccountNumber
	if paramNumber != "" {
		paramData["account_number"] = paramNumber
	}
	paramAccessKey := awsAccountEditAccessKeyID
	if paramAccessKey != "" {
		paramData["access_key_id"] = paramAccessKey
	}
	paramSecretKey := awsAccountEditSecretAccessKey
	if paramSecretKey != "" {
		paramData["secret_access_key"] = paramSecretKey
	}

	if len(paramData) <= 0 {
		cmd.Help()
		return errors.New("\nNone of the values to be changed has been specified.")
	}

	paramBytes, err := json.Marshal(paramData)
	if err != nil {
		return err
	}
	postBody := bytes.NewBuffer(paramBytes)

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/aws_accounts/{id}]
	httpRequest, err := client.createRequest(ctx, method, apiPath, currentProfile.APIKey, postBody)
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
