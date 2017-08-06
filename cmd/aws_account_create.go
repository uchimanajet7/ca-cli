package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "name" flag
var awsAccountCreateName string

// holds value of "account-number" flag
var awsAccountCreateAccountNumber string

// holds value of "access-key-id" flag
var awsAccountCreateAccessKeyID string

// holds value of "secret-access-key" flag
var awsAccountCreateSecretAccessKey string

var awsAccountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Register new CA AWS account.",
	Long:  `Use the "/aws_accounts" API to register new AWS account with Cloud Automator.`,
	RunE:  execAwsAccountCreate,
}

func init() {
	awsAccountCmd.AddCommand(awsAccountCreateCmd)

	awsAccountCreateCmd.Flags().StringVar(&awsAccountCreateName, "name", "", "Specify CA AWS account name. [required]")
	awsAccountCreateCmd.Flags().StringVar(&awsAccountCreateAccountNumber, "account-number", "", "Specify AWS account number. [required]")
	awsAccountCreateCmd.Flags().StringVar(&awsAccountCreateAccessKeyID, "access-key-id", "", "Specify AWS access key ID. [required]")
	awsAccountCreateCmd.Flags().StringVar(&awsAccountCreateSecretAccessKey, "secret-access-key", "", "Specify AWS secret access key. [required]")
}

func execAwsAccountCreate(cmd *cobra.Command, args []string) error {
	paramData := map[string]interface{}{}

	paramName := awsAccountCreateName
	if paramName == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--name] required CA AWS account name.")
	}
	paramData["name"] = paramName

	paramNumber := awsAccountCreateAccountNumber
	if paramNumber == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--account-number] required AWS account number.")
	}
	paramData["account_number"] = paramNumber

	paramAccessKey := awsAccountCreateAccessKeyID
	if paramAccessKey == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--access-key-id] required AWS access key ID.")
	}
	paramData["access_key_id"] = paramAccessKey

	paramSecretKey := awsAccountCreateSecretAccessKey
	if paramSecretKey == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--secret-access-key] required AWS secret access key.")
	}
	paramData["secret_access_key"] = paramSecretKey

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := "aws_accounts"
	method := "POST"
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

	// call API [/aws_accounts]
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
