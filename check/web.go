package check

import (
    "net/http"

    "time"
)

func Web(ip string) bool {
    httpClient := &http.Client{Timeout: time.Second}

    _, err := httpClient.Get(ip)

    return (err == nil)
}
