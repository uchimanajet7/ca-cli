package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// holds value of "name" flag
var postProcessCreateName string

// holds value of "service" flag
var postProcessCreateService string

// holds value of "parameters" flag
var postProcessCreateParameters string

var postProcessCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Register new CA post process",
	Long:  `Use the "/post_processes" API to register new post process with Cloud Automator.`,
	RunE:  execPostProcessCreate,
}

func init() {
	postProcessCmd.AddCommand(postProcessCreateCmd)

	postProcessCreateCmd.Flags().StringVar(&postProcessCreateName, "name", "", "Specify CA post process name. [required]")
	postProcessCreateCmd.Flags().StringVar(&postProcessCreateService, "service", "", "Specify CA post process service type. [required]")
	postProcessCreateCmd.Flags().StringVar(&postProcessCreateParameters, "parameters", "", "Specify CA post process setting values. [required]")
}

func execPostProcessCreate(cmd *cobra.Command, args []string) error {
	postBody, err := createPostProcessCreatePostBody()
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
	apiPath := "post_processes"
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

func createPostProcessCreatePostBody() (io.Reader, error) {
	params := map[string]interface{}{}

	paramName := postProcessCreateName
	if paramName == "" {
		return nil, errors.New("\nPlease specify flag [--name] required CA post process name.")
	}
	params["name"] = paramName

	paramService := postProcessCreateService
	if paramService == "" {
		return nil, errors.New("\nPlease specify flag [--service] required CA post process service type.")
	}
	params["service"] = paramService

	paramParameters := postProcessCreateParameters
	if paramParameters == "" {
		return nil, errors.New("\nPlease specify flag [--parameters] required CA post process setting values.")
	}
	parsedParam, err := createJobParseObjectParameter(paramParameters, "\nPlease specify flag [--parameters] required CA post process setting values.")
	if err != nil {
		return nil, err
	}
	params["parameters"] = *parsedParam

	paramBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(paramBytes), nil
}
