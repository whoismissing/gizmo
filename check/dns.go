package check

import (
	"net"
)

func Dns(ip string, record string) bool {
	foundIP, err := net.LookupIP(record)
	return (err == nil && ip == foundIP[0].String())
}
