package lib

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// callAPI creates a http.Request object, attaches headers to it and makes the
// requested api call.
func callAPI(method, url, postData string, headers map[string]string) (*http.Response, error) {
	var err error
	var req *http.Request

	if postData != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(postData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"API returned an error.\n\tMethod: %s\n\tURL: %s\n\tCode: %v\n\tStatus: %s\n\tRequest Body: %s\n\tResponse Body: %s",
			method, url, resp.StatusCode, resp.Status, postData, bodyBytes)
	}

	return resp, nil
}
