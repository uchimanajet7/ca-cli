package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "name" flag
var jobCreateName string

// holds value of "aws-account-id" flag
var jobCreateAwsAccountID string

// holds value of "rule-type" flag
var jobCreateRuleType string

// holds value of "rule-value" flag
var jobCreateRuleValue string

// holds value of "action-type" flag
var jobCreateActionType string

// holds value of "action-value" flag
var jobCreateActionValue string

// holds value of "completed-post-process-id" flag
var jobCreateCompletedPostProcessID string

// holds value of "failed-post-process-id" flag
var jobCreateFailedPostProcessID string

var jobCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Register new CA job",
	Long:  `Use the "/jobs" API to register new job with Cloud Automator.`,
	RunE:  execJobCreate,
}

func init() {
	jobCmd.AddCommand(jobCreateCmd)

	jobCreateCmd.Flags().StringVar(&jobCreateName, "name", "", "Specify CA job name. [required]")
	jobCreateCmd.Flags().StringVar(&jobCreateAwsAccountID, "aws-account-id", "", "Specify CA AWS account ID. [required]")
	jobCreateCmd.Flags().StringVar(&jobCreateRuleType, "rule-type", "", "Specify CA job trigger type. [required]")
	jobCreateCmd.Flags().StringVar(&jobCreateRuleValue, "rule-value", "", "Specify CA job trigger setting values. [required only \"--rule-type=cron/sqs\"]")
	jobCreateCmd.Flags().StringVar(&jobCreateActionType, "action-type", "", "Specify CA job action type. [required]")
	jobCreateCmd.Flags().StringVar(&jobCreateActionValue, "action-value", "", "Specify CA job action setting values. [required]")
	jobCreateCmd.Flags().StringVar(&jobCreateCompletedPostProcessID, "completed-post-process-id", "", "Specify array that contains post-processing IDs to be executed if the CA job succeeds.")
	jobCreateCmd.Flags().StringVar(&jobCreateFailedPostProcessID, "failed-post-process-id", "", "Specify array that contains post-processing IDs to be executed if the CA job faileds.")
}

func execJobCreate(cmd *cobra.Command, args []string) error {
	postBody, err := createJobCreatePostBody()
	if err != nil {
		cmd.Help()
		return err
	}

	// get profile data
	paramProfileName := rootPprofileName
	currentProfile, err := getProfile(paramProfileName)
	if err != nil {
		return err
	}

	// preparing to call API
	endURL := currentProfile.Endpoint
	apiPath := "jobs"
	method := "POST"

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/jobs]
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

func createJobCreatePostBody() (io.Reader, error) {
	params := map[string]interface{}{}

	paramName := jobCreateName
	if paramName == "" {
		return nil, errors.New("\nPlease specify flag [--name] required CA job name.")
	}
	params["name"] = paramName

	paramAwsAccountID := jobCreateAwsAccountID
	if paramAwsAccountID == "" {
		return nil, errors.New("\nPlease specify flag [--aws-account-id] required CA AWS account ID.")
	}
	params["aws_account_id"] = paramAwsAccountID

	paramRuleType := jobCreateRuleType
	if paramRuleType == "" {
		return nil, errors.New("\nPlease specify flag [--rule-type] required CA job trigger type.")
	}
	valid, required := isValidJobRuleType(paramRuleType)
	if !valid {
		return nil, errors.New("\nPlease specify flag [--rule-type] valid CA job trigger type.")
	}
	params["rule_type"] = paramRuleType

	paramRuleValue := jobCreateRuleValue
	if required {
		if paramRuleValue == "" {
			return nil, errors.New("\nPlease specify flag [--rule-value] required CA job trigger setting values.")
		}
		parsedParam, err := createJobParseObjectParameter(paramRuleValue, "\nPlease specify flag [--rule-value] required CA job trigger setting values.")
		if err != nil {
			return nil, err
		}
		params["rule_value"] = *parsedParam
	} else {
		params["rule_value"] = paramRuleValue
	}

	paramActionType := jobCreateActionType
	if paramActionType == "" {
		return nil, errors.New("\nPlease specify flag [--action-type] required CA job action type.")
	}
	params["action_type"] = paramActionType

	paramActionValue := jobCreateActionValue
	if paramActionValue == "" {
		return nil, errors.New("\nPlease specify flag [--action-value] required CA job action setting values.")
	}
	parsedParam, err := createJobParseObjectParameter(paramActionValue, "\nPlease specify flag [--action-value] required CA job action setting values.")
	if err != nil {
		return nil, err
	}
	params["action_value"] = *parsedParam

	// optional
	paramCompletedPostProcessID := jobCreateCompletedPostProcessID
	if paramCompletedPostProcessID != "" {
		params["completed_post_process_id"] = strings.Split(paramCompletedPostProcessID, ",")
	}

	// optional
	paramFailedPostProcessID := jobCreateFailedPostProcessID
	if paramFailedPostProcessID != "" {
		params["failed_post_process_id"] = strings.Split(paramFailedPostProcessID, ",")
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(paramBytes), nil
}
