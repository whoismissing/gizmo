package check

import (
	"net"
)

// Dns() makes a DNS standard request query for the user-specified A record
// and compares the response to the user-specified ip to check if they match
func Dns(ip string, record string) bool {
	foundIP, err := net.LookupIP(record)

	var status bool

	if err != nil {
		status = false
	}

	if ip == foundIP[0].String() {
		status = true
	}

	return status
}
