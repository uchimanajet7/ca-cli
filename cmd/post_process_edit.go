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
	"time"

	"github.com/spf13/cobra"
)

// holds value of "post-process-id" flag
var postProcessEditID string

// holds value of "name" flag
var postProcessEditName string

// holds value of "service" flag
var postProcessEditService string

// holds value of "parameters" flag
var postProcessEditParameters string

var postProcessEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Update CA post process",
	Long:  `Use the "/post_process/{id}" API to update specified ID post process with Cloud Automator.`,
	RunE:  execPostProcessEdit,
}

func init() {
	postProcessCmd.AddCommand(postProcessEditCmd)

	postProcessEditCmd.Flags().StringVar(&postProcessEditID, "post-process-id", "", "Specify CA post process ID. [required]")
	postProcessEditCmd.Flags().StringVar(&postProcessEditName, "name", "", "Specify CA post process name.")
	postProcessEditCmd.Flags().StringVar(&postProcessEditService, "service", "", "Specify CA post process service type.")
	postProcessEditCmd.Flags().StringVar(&postProcessEditParameters, "parameters", "", "Specify CA post process setting values.")
}

func execPostProcessEdit(cmd *cobra.Command, args []string) error {
	paramID := postProcessEditID
	if paramID == "" {
		cmd.Help()
		return errors.New("\nPlease specify flag [--post-process-id] required CA post process ID.")
	}

	postBody, err := createPostProcessEditPostBody()
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
	apiPath := path.Join("post_processes", paramID)
	method := "PATCH"

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

func createPostProcessEditPostBody() (io.Reader, error) {
	params := map[string]interface{}{}

	paramName := postProcessEditName
	if paramName != "" {
		params["name"] = paramName
	}

	paramService := postProcessEditService
	if paramService != "" {
		params["service"] = paramService
	}

	paramParameters := postProcessEditParameters
	if paramParameters != "" {
		parsedParam, err := createJobParseObjectParameter(paramParameters, "\nPlease specify flag [--parameters] required CA post process setting values.")
		if err != nil {
			return nil, err
		}
		params["parameters"] = *parsedParam
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
