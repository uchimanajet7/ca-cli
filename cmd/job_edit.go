package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "job-id" flag
var jobEditID string

// holds value of "name" flag
var jobEditName string

// holds value of "aws-account-id" flag
var jobEditAwsAccountID string

// holds value of "rule-value" flag
var jobEditRuleValue string

// holds value of "action-value" flag
var jobEditActionValue string

// holds value of "completed-post-process-id" flag
var jobEditCompletedPostProcessID string

// holds value of "failed-post-process-id" flag
var jobEditFailedPostProcessID string

var jobEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Update CA job",
	Long:  `Use the "/jobs{id}" API to update specified ID job with Cloud Automator.`,
	RunE:  execJobEdit,
}

func init() {
	jobCmd.AddCommand(jobEditCmd)

	jobEditCmd.Flags().StringVar(&jobEditID, "job-id", "", "Specify CA job ID. [required]")
	jobEditCmd.Flags().StringVar(&jobEditName, "name", "", "Specify CA job name.")
	jobEditCmd.Flags().StringVar(&jobEditAwsAccountID, "aws-account-id", "", "Specify CA AWS account ID.")
	jobEditCmd.Flags().StringVar(&jobEditRuleValue, "rule-value", "", "Specify CA job trigger setting values.")
	jobEditCmd.Flags().StringVar(&jobEditActionValue, "action-value", "", "Specify CA job action setting values.")
	jobEditCmd.Flags().StringVar(&jobEditCompletedPostProcessID, "completed-post-process-id", "", "Specify array that contains post-processing IDs to be executed if the CA job succeeds.")
	jobEditCmd.Flags().StringVar(&jobEditFailedPostProcessID, "failed-post-process-id", "", "Specify array that contains post-processing IDs to be executed if the CA job faileds.")
}

func execJobEdit(cmd *cobra.Command, args []string) error {
	paramID := jobEditID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--job-id] required CA job ID.")
	}

	postBody, err := createJobEditPostBody()
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
	apiPath := path.Join("jobs", paramID)
	method := "PATCH"

	// create http client
	httpClient := &http.Client{}
	client, err := createClient(endURL, httpClient)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// call API [/jobs/{id}]
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

func createJobEditPostBody() (io.Reader, error) {
	params := map[string]interface{}{}

	paramName := jobEditName
	if paramName != "" {
		params["name"] = paramName
	}

	paramAwsAccountID := jobEditAwsAccountID
	if paramAwsAccountID != "" {
		params["aws_account_id"] = paramAwsAccountID
	}

	paramRuleValue := jobEditRuleValue
	if paramRuleValue != "" {
		parsedParam, err := createJobParseObjectParameter(paramRuleValue, "")
		if err == nil {
			params["rule_value"] = *parsedParam
		}
	}

	paramActionValue := jobEditActionValue
	if paramActionValue != "" {
		parsedParam, err := createJobParseObjectParameter(paramActionValue, "")
		if err == nil {
			params["action_value"] = *parsedParam
		}
	}

	paramCompletedPostProcessID := jobEditCompletedPostProcessID
	if paramCompletedPostProcessID != "" {
		params["completed_post_process_id"] = strings.Split(paramCompletedPostProcessID, ",")
	}

	paramFailedPostProcessID := jobEditFailedPostProcessID
	if paramFailedPostProcessID != "" {
		params["failed_post_process_id"] = strings.Split(paramFailedPostProcessID, ",")
	}

	if len(params) <= 0 {
		return nil, errors.New("\nNone of the values to be changed has been specified.")
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(paramBytes), nil
}
