package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

// Client use HTTP access
type Client struct {
	EndpointURL *url.URL
	HTTPClient  *http.Client
}

func createClient(endpointURL string, httpClient *http.Client) (*Client, error) {
	parsedURL, err := url.Parse(endpointURL)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse url: %s", endpointURL)
	}

	client := &Client{
		EndpointURL: parsedURL,
		HTTPClient:  httpClient,
	}
	return client, nil
}

func (client *Client) createRequest(ctx context.Context, method string, subPath string, apiKey string, body io.Reader) (*http.Request, error) {
	endpointURL := *client.EndpointURL
	endpointURL.Path = path.Join(client.EndpointURL.Path, subPath)

	req, err := http.NewRequest(method, endpointURL.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ca-cli/v000")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func getBodyString(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// check http error
	if resp.StatusCode >= http.StatusBadRequest {
		// to-dp: JSON format may be more convenient?
		return "", fmt.Errorf("\nResponse Code: %d\nResponse Body: %s\n", resp.StatusCode, string(b))
	}

	// body data is none
	if len(b) <= 0 {
		b = []byte(fmt.Sprintf("{\"response_code\":\"%d\"}", resp.StatusCode))
	}
	// indent 4 space
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, b, "", "    "); err != nil {
		return "", err
	}

	return prettyJSON.String(), nil
}

func dumpHTTPRequest(req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(dump))
}

func dumpHTTPResponse(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(dump))
}
