package lib

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// CallAPI creates a http.Request object, attaches headers to it and makes the
// requested api call.
func CallAPI(method, url, postData string, headers map[string]string) *http.Response {
	var err error
	var req *http.Request

	if postData != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(postData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if resp.StatusCode >= 300 {
		fmt.Println(fmt.Sprintf(
			"API returned an error. \n\tMethod: %s, \n\tURL: %s, \n\tCode: %v, \n\tStatus: %s \n\tBody: %s",
			method, url, resp.StatusCode, resp.Status, postData))
		os.Exit(1)
	}

	return resp
}
