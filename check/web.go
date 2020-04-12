package check

import (
	"net/http"
	"time"
)

// Web() makes a GET request to a url in the formats:
// http://exampleone.com, https://exampletwo.com
// http://192.168.1.1
func Web(url string) bool {
	// TODO: make timeout adjustable
	httpClient := &http.Client{Timeout: time.Second}
	_, err := httpClient.Get(url)

	var status bool
	if err != nil {
		status = false
	} else {
		status = true
	}

	return status
}
